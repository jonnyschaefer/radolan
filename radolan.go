// Package radolan parses the DWD RADOLAN / RADVOR radar composite format. This data
// is available at the Global Basic Dataset (http://www.dwd.de/DE/leistungen/gds/gds.html).
// The obtained results can be processed and visualized with additional functions.
//
// Currently the national grid [1][4] and the extended european grid [5] are supported.
// Tested input products are PG, FZ, SF, RW, RX and EX. Those can be considered working with
// sufficient accuracy.
//
// In cases, where the publicly available format specification is unprecise or contradictory,
// reverse engineering was used to obtain reasonable approaches.
//	Used references:
//
//	[1] https://www.dwd.de/DE/leistungen/radolan/radolan_info/radolan_radvor_op_komposit_format_pdf.pdf
//	[2] https://www.dwd.de/DE/leistungen/gds/weiterfuehrende_informationen.zip
//	[3]  - legend_radar_products_fz_forecast.pdf
//	[4]  - legend_radar_products_pg_coordinates.pdf
//	[5]  - legend_radar_products_radolan_rw_sf.pdf
//	[6] https://www.dwd.de/DE/leistungen/radarniederschlag/rn_info/download_niederschlagsbestimmung.pdf
package radolan

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

// Radolan radar data is provided in so called composite formats. Each composite is a combined
// image consisting of mulitiple radar sweeps spread over the composite area.
// The composite c has a an internal resolution of c.Dx (horizontal) * c.Dy (vertical) records
// covering a real surface of c.Dx * c.Rx * c.Dy * c.Dy square kilometers.
// The pixel value at the position (x, y) is represented by c.Data[ y ][ x ] and is stored as
// raw rvp-6 value (NaN if the no-data flag is set). This rvp-6 value is used differently
// depending on the product type:
//
//	Product label    ||   raw value        "live" cloud reflectivity       "live" rainfall rate
//	-----------------||   +-------+              +-----+                          +------+
//	(PG), FZ, ...    ||   | rvp-6 |---ToDBZ()--->| dBZ |---PrecipitationRate()--->| mm/h |
//	 RX, EX          ||   +---+---+              +-----+                          +------+
//	-----------------|| - - - | - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//	 RW, SF,  ...    ||       |       +------+
//	                 ||       +----- >| mm/h |
//	                 ||               | mm/d |
//	-----------------||               +------+
//	                            aggregated precipitation
//
// The cloud reflectivity factor Z is stored in its logarithmic representation dBZ:
//	dBZ = 10 * log(Z)
// Real world geographical coordinates (latitude, longitude) can be projected into the
// coordinate system of the composite by using the translation method:
//	x, y := c.Translate(52.51861, 13.40833)	// Berlin (lat, lon)
//
//	rvp := c.At(int(x), int(y))				// Raw value (rvp-6)
//	dbz := rvp.ToDBZ()					// Cloud reflectivity (dBZ)
//	rat := dbz.PrecipitationRate(radolan.Doelling98)	// Rainfall rate (mm/h) using Doelling98 as Z-R relationship
//
//	fmt.Println("Rainfall in Berlin [mm/h]:", rat)
//
type Composite struct {
	Product string // composite product label

	CaptureTime  time.Time     // time of source data capture used for forcasting
	ForecastTime time.Time     // data represents conditions predicted for this time
	Interval     time.Duration // time duration until next forecast

	Data [][]RVP6 // rvp-6 data for each point [y][x]

	Dx int // data width
	Dy int // data height

	Rx float64 // horizontal resolution in km/px
	Ry float64 // vertical resolution in km/px

	dataLength int // length of binary section in bytes

	precision int   // multiplicator 10^precision for each raw value
	level     []DBZ // maps data value to corresponding dBZ value in runlength based formats

	offx float64 // horizontal projection offset
	offy float64 // vertical projection offset
}

// NewComposite reads binary data from rd and parses the composite.
func NewComposite(rd io.Reader) (comp *Composite, err error) {
	reader := bufio.NewReader(rd)
	comp = &Composite{}

	err = comp.parseHeader(reader)
	if err != nil {
		return
	}

	err = comp.parseData(reader)
	if err != nil {
		return
	}

	comp.calibrateProjection()

	return
}

// NewComposites reads tar gz data from rd and returns the parsed composites sorted by
// ForecastTime in ascending order.
func NewComposites(rd io.Reader) ([]*Composite, error) {
	gzipReader, err := gzip.NewReader(rd)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	var cs []*Composite
	for {
		_, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		c, err := NewComposite(tarReader)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}

	// sort composites in chronological order
	sort.Slice(cs, func(i, j int) bool { return cs[i].ForecastTime.Before(cs[j].ForecastTime) })
	return cs, nil
}

// NewDummy creates a blank dummy composite with the given product label and dimensions. It can
// be used for generic coordinate translation.
func NewDummy(product string, dx, dy int) (comp *Composite) {
	comp = &Composite{Product: product, Dx: dx, Dy: dy}
	comp.calibrateProjection()
	return
}

// At is shorthand for c.Data[y][x] and returns the radar video processor value (rvp-6) at
// the given point. NaN is returned, if no data is available or the requested point is located
// outside the scanned area.
func (c *Composite) At(x, y int) RVP6 {
	if x < 0 || y < 0 || x >= c.Dx || y >= c.Dy {
		return RVP6(math.NaN())
	}
	return c.Data[y][x]
}

// newError returns an error indicating the failed function and reason
func newError(function, reason string) error {
	return fmt.Errorf("radolan.%s: %s", function, reason)
}

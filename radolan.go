// Package radolan parses the DWD RADOLAN / RADVOR radar composite format. This data
// is available at the Open Data Portal (https://www.dwd.de/DE/leistungen/opendata/opendata.html).
// The obtained results can be processed and visualized with additional functions.
//
// Tested input products and grids:
//
//	| Product | Grid              | Description             |
//	| ------- | ----------------- | ----------------------- |
//	| EX      | middle-european   | reflectivity            |
//	| FX      | national          | nowcast reflectivity    |
//	| FZ      | national          | nowcast reflectivity    |
//	| PE      | local             | echo top                |
//	| PF      | local             | reflectivity            |
//	| PG      | national picture  | reflectivity            |
//	| PR      | local             | doppler radial velocity |
//	| PX      | local             | reflectivity            |
//	| PZ      | local             | 3D reflectivity CAPPI   |
//	| RW      | national          | hourly accumulated      |
//	| RX      | national          | reflectivity            |
//	| SF      | national          | daily accumulated       |
//	| WX      | extended national | reflectivity            |
//
// Those can be considered working with sufficient accuracy.
//
// In cases, where the publicly available format specification is unprecise or contradictory,
// reverse engineering was used to obtain reasonable approaches.
//	Used references:
//
//	[1] https://www.dwd.de/DE/leistungen/radolan/radolan_info/radolan_radvor_op_komposit_format_pdf.pdf
//	[2] https://www.dwd.de/DE/leistungen/gds/weiterfuehrende_informationen.zip
//	[3] https://www.dwd.de/DE/leistungen/radarprodukte/formatbeschreibung_fxdaten.pdf
//	[4] https://www.dwd.de/DE/leistungen/opendata/help/radar/radar_pg_coordinates_pdf.pdf
//	[5] https://www.dwd.de/DE/leistungen/radarniederschlag/rn_info/download_niederschlagsbestimmung.pdf
//	[6] https://www.dwd.de/DE/leistungen/radarprodukte/formatbeschreibung_wndaten.pdf
//	[7] hex editor and much reverse engineering
package radolan

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"fmt"
	"io"
	"sort"
	"time"
)

// Radolan radar data is provided as single local sweeps or so called composite
// formats. Each composite is a combined image consisting of mulitiple radar
// sweeps spread over the composite area.
// The 2D composite c has a an internal resolution of c.Dx (horizontal) * c.Dy
// (vertical) records covering a real surface of c.Dx * c.Rx * c.Dy * c.Dy
// square kilometers.
// The pixel value at the position (x, y) is represented by
// c.Data[ y ][ x ] and is stored as raw float value (NaN if the no-data flag
// is set). Some 3D radar products feature multiple layers in which the voxel
// at position (x, y, z) is accessible by c.DataZ[ z ][ y ][ x ].
//
// The data value is used differently depending on the product type:
// (also consult the DataUnit field of the Composite)
//
//	Product label            | values represent         | unit
//	-------------------------+--------------------------+------------------------
//	 PG, PC, PX*, ...        | cloud reflectivity       | dBZ
//	 RX, WX, EX, FZ, FX, ... | cloud reflectivity	    | dBZ
//	 RW, SF,  ...            | aggregated precipitation | mm/interval
//	 PR*, ...                | doppler radial velocity  | m/s
//
// The cloud reflectivity (in dBZ) can be converted to rainfall rate (in mm/h)
// via PrecipitationRate().
//
// The cloud reflectivity factor Z is stored in its logarithmic representation dBZ:
//	dBZ = 10 * log(Z)
// Real world geographical coordinates (latitude, longitude) can be projected into the
// coordinate system of the composite by using the projection method:
//	// if c.HasProjection
//	x, y := c.Project(52.51861, 13.40833)	// Berlin (lat, lon)
//
//	dbz := c.At(int(x), int(y))					// Raw value is Cloud reflectivity (dBZ)
//	rat := radolan.PrecipitationRate(radolan.Doelling98, dbz)	// Rainfall rate (mm/h) using Doelling98 as Z-R relationship
//
//	fmt.Println("Rainfall in Berlin [mm/h]:", rat)
//
type Composite struct {
	Product string // composite product label

	CaptureTime  time.Time     // time of source data capture used for forcasting
	ForecastTime time.Time     // data represents conditions predicted for this time
	Interval     time.Duration // time duration until next forecast

	DataUnit Unit

	PlainData [][]float32 // data for parsed plain data element [y][x]
	Px        int         // plain data width
	Py        int         // plain data height

	DataZ [][][]float32 // data for each voxel [z][y][x] (composites use only one z-layer)
	Data  [][]float32   // data for each pixel [y][x] at layer 0 (alias for DataZ[0][x][y])

	Dx int // data width
	Dy int // data height
	Dz int // data layer

	Rx float64 // horizontal resolution in km/px
	Ry float64 // vertical resolution in km/px

	HasProjection bool // coordinate projection available

	Format int // Version Format

	dataLength int // length of binary section in bytes

	precision int       // multiplicator 10^precision for each raw value
	level     []float32 // maps data value to corresponding index value in runlength based formats

	offx float64 // horizontal projection offset
	offy float64 // vertical projection offset

	proj_wgs84 *projection
}

// ErrUnknownUnit indicates that the unit of the radar data is not defined in
// the catalog. The data values can be incorrect due to unit dependent
// conversions during parsing. Be careful when further processing the
// composite.
var ErrUnknownUnit = newError("NewComposite", "data unit not defined in catalog. data values can be incorrect")

// NewComposite reads binary data from rd and parses the composite.  An error
// is returned on failure. When ErrUnknownUnit is returned, the data values can
// be incorrect due to unit dependent conversions during parsing. In this case
// be careful when further processing the composite.
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
	comp.arrangeData()

	comp.calibrateProjection()

	if comp.DataUnit == Unit_unknown {
		err = ErrUnknownUnit
	}

	return
}

// NewComposites reads .tar.bz2 data from rd and returns the parsed composites sorted by
// ForecastTime in ascending order.
func NewComposites(rd io.Reader) ([]*Composite, error) {
	bzipReader := bzip2.NewReader(rd)

	tarReader := tar.NewReader(bzipReader)

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

// NewDummy creates a blank dummy composite with the given product label, format version, and dimensions. It can
// be used for generic coordinate projection.
func NewDummy(product string, format, dx, dy int) (comp *Composite) {
	comp = &Composite{Product: product, Format: format, Dx: dx, Dy: dy}
	comp.calibrateProjection()
	return
}

// At is shorthand for c.Data[y][x] and returns the radar video processor value
// at the given point. NaN is returned, if no data is available or the
// requested point is located outside the scanned area.
func (c *Composite) At(x, y int) float32 {
	return c.AtZ(x, y, 0)
}

// AtZ is shorthand for c.DataZ[z][y][x] and returns the radar video processor
// value at the given point. NaN is returned, if no data is available or the
// requested point is located outside the scanned volume.
func (c *Composite) AtZ(x, y, z int) float32 {
	if x < 0 || y < 0 || z < 0 || x >= c.Dx || y >= c.Dy || z >= c.Dz {
		return NaN
	}
	return c.DataZ[z][y][x]
}

// newError returns an error indicating the failed function and reason
func newError(function, reason string) error {
	return fmt.Errorf("radolan.%s: %s", function, reason)
}

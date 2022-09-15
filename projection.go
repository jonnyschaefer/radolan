package radolan

import (
	"math"
)

// corner coordinates for projection are defined by the grid type
type grid int

const (
	unknownGrid          grid = iota
	nationalGrid              // resolution: 900km * 900km
	nationalPictureGrid       // resolution: 920km * 920km
	extendedNationalGrid      // resolution: 900km * 1100km
	DE1200Grid                // resolution: 1100km * 1200km
	middleEuropeanGrid        // resolution: 1400km * 1500km
)

// minRes repeatedly bisects the given edges until no further step is possible
// for at least one edge. The resulting dimensions are the returned.
func minRes(dx, dy int) (rdx int, rdy int) {
	rdx, rdy = dx, dy

	if rdx == 0 || rdy == 0 {
		return
	}
	for rdx&1 == 0 && rdy&1 == 0 {
		rdx >>= 1
		rdy >>= 1
	}
	return
}

// errNoProjection means that the projection grid could not be identified.
var errNoProjection = newError("cornerPoints", "warning: unable to identify grid")

// detectGrid identifies the used projection grid based on the composite dimensions
func (c *Composite) detectGrid() grid {
	dx, dy := minRes(c.Dx, c.Dy)

	if mx, my := minRes(900, 900); dx == mx && dy == my {
		return nationalGrid
	}
	if mx, my := minRes(920, 920); dx == mx && dy == my {
		return nationalPictureGrid
	}
	if mx, my := minRes(900, 1100); dx == mx && dy == my {
		return extendedNationalGrid
	}
	if mx, my := minRes(1100, 1200); dx == mx && dy == my {
		return DE1200Grid
	}
	if mx, my := minRes(1400, 1500); dx == mx && dy == my {
		return middleEuropeanGrid
	}
	return unknownGrid
}

// cornerPoints returns corner coordinates of the national, extended or
// middle-european grid based on the dimensions of the composite. The used
// values are described in [1], [4].  If an error is returned,
// projection methods will not work.
func (c *Composite) cornerPoints(detectedGrid grid) (originTop, originLeft, edgeBottom, edgeRight float64, err error) {
	switch detectedGrid {
	case nationalGrid: // described in [1]
		originTop, originLeft = 54.5877, 02.0715 // N, E
		edgeBottom, edgeRight = 47.0705, 14.6209 // N, E
	case nationalPictureGrid: // (pg) described in [4]
		originTop, originLeft = 54.66218275, 1.900684377 // N, E
		edgeBottom, edgeRight = 46.98044293, 14.73300934 // N, E
	case extendedNationalGrid: // described in [1]
		originTop, originLeft = 55.5482, 03.0889 // N, E
		edgeBottom, edgeRight = 46.1827, 15.4801 // N, E
	case DE1200Grid: // described in [6]
		if c.Format >= 5 {
			originTop, originLeft = 55.86208711, 1.463301510 // N, E
			edgeBottom, edgeRight = 45.68460578, 16.58086935 // N, E
		} else {
			originTop, originLeft = 55.86584289, 1.435612143 // N, E
			edgeBottom, edgeRight = 45.68358331, 16.60186543 // N, E
		}
	case middleEuropeanGrid: // described in [1]
		originTop, originLeft = 56.5423, -0.8654 // N, E
		edgeBottom, edgeRight = 43.8736, 18.2536 // N, E
	default:
		err = errNoProjection
	}
	return
}

// calibrateProjection initializes fields that are necessary for coordinate transformation
func (c *Composite) calibrateProjection() {
	c.Rx = math.NaN()
	c.Ry = math.NaN()
	c.offx = math.NaN()
	c.offy = math.NaN()

	// get corner points
	detectedGrid := c.detectGrid()

	originTop, originLeft, edgeBottom, edgeRight, err := c.cornerPoints(detectedGrid)
	if err != nil {
		return
	}

	if c.Format >= 5 {
		if detectedGrid == DE1200Grid {
			c.proj_wgs84 = proj_DE1200_WGS84
		} else {
			return
		}
	}

	// found matching projection rule
	c.HasProjection = true

	// set resolution to 1 km for calibration
	c.Rx = 1.0
	c.Ry = 1.0
	c.offx = 0.0
	c.offy = 0.0

	// calibrate offset correction
	c.offx, c.offy = c.Project(originTop, originLeft)

	// calibrate scaling
	resx, resy := c.Project(edgeBottom, edgeRight)
	c.Rx = (resx) / float64(c.Dx)
	c.Ry = (resy) / float64(c.Dy)
}


// Project transforms geographical coordinates (latitude north, longitude east) to the
// according data indices in the coordinate system of the composite.
// NaN is returned when no projection is available. Procedures adapted from [1] and [6].
func (c *Composite) Project(north, east float64) (x, y float64) {
	if !c.HasProjection {
		x, y = math.NaN(), math.NaN()
		return
	}

	if c.proj_wgs84 != nil {
		return c.projectWGS84(north, east)
	}
	return c.projectSphere(north, east)
}

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
	middleEuropeanGrid        // resolution: 1400km * 1500km
)

// values described in [1]
const (
	earthRadius = 6370.04 // km

	junctionNorth = 60.0 // N
	junctionEast  = 10.0 // E
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
	if mx, my := minRes(1400, 1500); dx == mx && dy == my {
		return middleEuropeanGrid
	}
	return unknownGrid
}

// cornerPoints returns corner coordinates of the national, extended or
// middle-european grid based on the dimensions of the composite. The used
// values are described in [1], [4].  If an error is returned,
// translation methods will not work.
func (c *Composite) cornerPoints() (originTop, originLeft, edgeBottom, edgeRight float64, err error) {
	switch c.detectGrid() {
	case nationalGrid: // described in [1]
		originTop, originLeft = 54.5877, 02.0715 // N, E
		edgeBottom, edgeRight = 47.0705, 14.6209 // N, E
	case nationalPictureGrid: // (pg) described in [4]
		originTop, originLeft = 54.66218275, 1.900684377 // N, E
		edgeBottom, edgeRight = 46.98044293, 14.73300934 // N, E
	case extendedNationalGrid: // described in [1]
		originTop, originLeft = 55.5482, 03.0889 // N, E
		edgeBottom, edgeRight = 46.1827, 15.4801 // N, E
	case middleEuropeanGrid: // described in [1]
		originTop, originLeft = 56.5423, -0.8654 // N, E
		edgeBottom, edgeRight = 43.8736, 18.2536 // N, E
	default:
		err = errNoProjection
	}
	return
}

// calibrateProjection initializes fields that are necessary for coordinate translation
func (c *Composite) calibrateProjection() {
	// calibration is only neccessary for sperical projection
	if c.wgs84 {
		return
	}

	// get corner points
	originTop, originLeft, edgeBottom, edgeRight, err := c.cornerPoints()
	if err != nil {
		nan := math.NaN()
		c.Rx, c.Ry, c.offx, c.offy = nan, nan, nan, nan
		return
	}

	// found matching projection rule
	c.HasProjection = true

	// set resolution to 1 km for calibration
	c.Rx, c.Ry = 1.0, 1.0

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
	return c.projectSphere(north, east)
}

func (c *Composite) projectSphere(north, east float64) (x, y float64) {
	if !c.HasProjection {
		x, y = math.NaN(), math.NaN()
		return
	}

	rad := func(deg float64) float64 {
		return deg * (math.Pi / 180.0)
	}

	lamda0, phi0 := rad(junctionEast), rad(junctionNorth)
	lamda, phi := rad(east), rad(north)

	m := (1.0 + math.Sin(phi0)) / (1.0 + math.Sin(phi))
	x = (earthRadius * m * math.Cos(phi) * math.Sin(lamda-lamda0))
	y = (earthRadius * m * math.Cos(phi) * math.Cos(lamda-lamda0))

	// offset correction
	x -= c.offx
	y -= c.offy

	// scaling
	x /= c.Rx
	y /= c.Ry

	return
}

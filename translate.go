package radolan

import (
	"math"
)

// values described in [1]
const (
	earthRadius = 6370.04 // km

	junctionNorth = 60.0 // N
	junctionEast  = 10.0 // E
)

// isScaled returns true if the composite can be scaled to the given dimensions
// while also maintaining the aspect ratio.
func (c *Composite) isScaled(dx, dy int) bool {
	epsilon := 0.00000001
	return math.Abs(1-float64(c.Dy*dx)/float64(c.Dx*dy)) < epsilon
}

// errNoProjection means that the projection grid could not be identified.
var errNoProjection = newError("cornerPoints", "warning: unable to identify grid")

// cornerPoints returns corner coordinates of the national, extended or middle-european grid
// based on the product label or resolution of the composite. The used values are
// described in [1], [4] and [5].
// If an error is returned, translation methods will not work.
func (c *Composite) cornerPoints() (originTop, originLeft, edgeBottom, edgeRight float64, err error) {
	// national grid (pg) values described in [4]
	if c.Product == "PG" {
		originTop, originLeft = 54.6547, 01.9178 // N, E
		edgeBottom, edgeRight = 46.9894, 14.7218 // N, E
		return
	}

	// national grid values described in [1]
	if c.isScaled(900, 900) {
		originTop, originLeft = 54.5877, 02.0715 // N, E
		edgeBottom, edgeRight = 47.0705, 14.6209 // N, E
		return
	}

	// extended national grid described in [5]
	if c.isScaled(900, 1100) {
		originTop, originLeft = 55.5482, 03.0889 // N, E
		edgeBottom, edgeRight = 46.1827, 15.4801 // N, E
		return
	}

	// middle european grid described in [5]
	if c.isScaled(1400, 1500) {
		originTop, originLeft = 56.5423, -0.8654 // N, E
		edgeBottom, edgeRight = 43.8736, 18.2536 // N, E
		return
	}

	err = errNoProjection
	return
}

// calibrateProjection initializes fields that are necessary for coordinate translation
func (c *Composite) calibrateProjection() {
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
	c.offx, c.offy = c.Translate(originTop, originLeft)

	// calibrate scaling
	resx, resy := c.Translate(edgeBottom, edgeRight)
	c.Rx = (resx) / float64(c.Dx)
	c.Ry = (resy) / float64(c.Dy)
}

// Translate translates geographical coordinates (latitude north, longitude east) to the
// according data indices in the coordinate system of the composite.
// NaN is returned when no projection is available. Procedures adapted from [1].
func (c *Composite) Translate(north, east float64) (x, y float64) {
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

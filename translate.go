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

// cornerPoints returns corner coordinates of the national or extended european grid
// based on the product label or resolution of the composite. The used values are
// described in [1], [4] and [5].
func (c *Composite) cornerPoints() (originTop, originLeft, edgeBottom, edgeRight float64) {
	// national grid (pg) values described in [4]
	if c.Product == "PG" {
		originTop, originLeft = 54.6547, 01.9178 // N, E
		edgeBottom, edgeRight = 46.9894, 14.7218 // N, E
		return
	}

	// national grid values described in [1]
	if c.Dx == c.Dy {
		originTop, originLeft = 54.5877, 02.0715 // N, E
		edgeBottom, edgeRight = 47.0705, 14.6209 // N, E
		return
	}

	// extended european grid described in [5]
	originTop, originLeft = 56.5423, -0.8654 // N, E
	edgeBottom, edgeRight = 43.8736, 18.2536 // N, E
	return
}

// calibrateProjection initializes fields that are necessary for coordinate translation
func (c *Composite) calibrateProjection() {
	// get corner points
	originTop, originLeft, edgeBottom, edgeRight := c.cornerPoints()

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
// according data indices in the coordinate system of the composite. Procedures adapted from [1].
func (c *Composite) Translate(north, east float64) (x, y float64) {
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

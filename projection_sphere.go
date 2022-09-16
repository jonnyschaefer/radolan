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

func (c *Composite) projectSphere(north, east float64) (x, y float64) {
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

package radolan

import (
	"math"
)


type proj struct {
	lon_0 float64
	ecc float64
	k_0 float64
	x_0 float64
	y_0 float64

	postOffsetX float64
	postOffsetY float64
	postScaleX float64
	postScaleY float64
}

const (
	degToRad = 2 * math.Pi / 360
)

// DE1200 WGS84
// +proj=stere +lat_0=90 +lat_ts=60 +lon_0=10 +a=6378137 +b=6356752.3142451802 +no_defs +x_0=543196.83521776402 +y_0=3622588.861931001
var proj_DE1200_WGS84 = &proj{
	lon_0: 10 * degToRad,
	ecc: 0.08181919084262032,
	k_0: 11862667.042661695,
	x_0: 543196.83521776402,
	y_0: 3622588.861931001,

	postOffsetX: 0.5,
	postOffsetY: 0.5,
	postScaleX: 1000,
	postScaleY: -1000,
}

func (p *proj) project(lat, lon float64) (x, y float64) {
	lat *= degToRad
	lon *= degToRad

	sinLat := math.Sin(lat)

	s := p.k_0 * math.Tan(0.5 * (math.Pi / 2 -  lat)) / math.Pow(((1 - p.ecc * sinLat) / (1 + p.ecc * sinLat)), 0.5 * p.ecc)

	x = p.x_0 + (s * math.Sin(lon - p.lon_0))
	y = p.y_0 - (s * math.Cos(lon - p.lon_0))

	x = (x/p.postScaleX) + p.postOffsetX
	y = (y/p.postScaleY) + p.postOffsetY
	return
}

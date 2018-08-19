package radolan

import (
	"math"
)

var NaN = float32(math.NaN())

func IsNaN(f float32) (is bool) {
	return f != f
}

// Z-R relationship
type ZR struct {
	// intermediate caching
	c1 float64 // 10*b
	c2 float64 // a^(-1/b)
	c3 float64 // 10^(1/(10*b))
	c4 float64 // 10 * log10(a)
}

// Common Z-R relationships
var (
	Aniol80          = NewZR(256, 1.42) // operational use in germany, described in [5]
	Doelling98       = NewZR(316, 1.50) // operational use in switzerland
	JossWaldvogel70  = NewZR(300, 1.50)
	MarshallPalmer55 = NewZR(200, 1.60) // operational use in austria
)

// New Z-R returns a Z-R relationship mathematically expressed as Z = a * R^b
func NewZR(A, B float64) ZR {
	c1 := 10.0 * B
	c2 := math.Pow(A, -1.0/B)
	c3 := math.Pow(10.0, 1/c1)
	c4 := 10.0 * math.Log10(A)

	return ZR{c1, c2, c3, c4}
}

// PrecipitationRate returns the estimated precipitation rate in mm/h for the given
// reflectivity factor and Z-R relationship.
func PrecipitationRate(relation ZR, dBZ float32) (rate float64) {
	return relation.c2 * math.Pow(relation.c3, float64(dBZ))
}

// Reflectivity returns the estimated reflectivity factor for the given precipitation
// rate (mm/h) and Z-R relationship.
func Reflectivity(relation ZR, rate float64) (dBZ float32) {
	return float32(relation.c4 + relation.c1*math.Log10(rate))
}

// toDBZ converts the given radar video processor values (rvp-6) to radar reflectivity
// factors in decibel relative to Z (dBZ).
func toDBZ(rvp6 float32) (dBZ float32) {
	return rvp6/2.0 - 32.5
}

// toRVP6 converts the given radar reflectivity factors (dBZ) to radar video processor
// values (rvp-6).
func toRVP6(dBZ float32) float32 {
	return (dBZ + 32.5) * 2
}

// rvp6Raw converts the raw value to radar video processor values (rvp-6) by applying the
// products precision field.
func (c *Composite) rvp6Raw(value int) float32 {
	return float32(value) * float32(math.Pow10(c.precision))
}

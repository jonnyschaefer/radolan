package radolan

import (
	"math"
)

var NaN = float32(math.NaN())

func IsNaN(f float32) (is bool) {
	return f != f
}

// Z-R relationship mathematically expressed as Z = a * R^b
type ZR struct {
	A float64
	B float64
}

// Common Z-R relationships
var (
	Aniol80          = ZR{256, 1.42} // operational use in germany, described in [6]
	Doelling98       = ZR{316, 1.50} // operational use in switzerland
	JossWaldvogel70  = ZR{300, 1.50}
	MarshallPalmer55 = ZR{200, 1.60} // operational use in austria
)

// PrecipitationRate returns the estimated precipitation rate in mm/h for the given
// reflectivity factor and Z-R relationship.
func PrecipitationRate(relation ZR, dBZ float32) (rate float64) {
	return math.Pow(math.Pow(10, float64(dBZ)/10)/relation.A, 1/relation.B)
}

// Reflectivity returns the estimated reflectivity factor for the given precipitation
// rate (mm/h) and Z-R relationship.
func Reflectivity(relation ZR, rate float64) (dBZ float32) {
	return float32(10 * math.Log10(relation.A*math.Pow(float64(rate), relation.B)))
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

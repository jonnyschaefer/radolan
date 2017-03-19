package radolan

import (
	"math"
)

// Radar reflectivity factor Z in logarithmic representation dBZ: dBZ = 10 * log(Z)
type DBZ float64

// Raw radar video processor value.
type RVP6 float64

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
func (z DBZ) PrecipitationRate(relation ZR) float64 {
	return math.Pow(math.Pow(10, float64(z)/10)/relation.A, 1/relation.B)
}

// Reflectivity returns the estimated reflectivity factor for the given precipitation
// rate (mm/h) and Z-R relationship.
func Reflectivity(rate float64, relation ZR) DBZ {
	return DBZ(10 * math.Log10(relation.A*math.Pow(rate, relation.B)))
}

// ToDBZ converts the given radar video processor values (rvp-6) to radar reflectivity
// factors in decibel relative to Z (dBZ).
func (r RVP6) ToDBZ() DBZ {
	return DBZ(r/2.0 - 32.5)
}

// ToRVP6 converts the given radar reflectivity factors (dBZ) to radar video processor
// values (rvp-6).
func (z DBZ) ToRVP6() RVP6 {
	return RVP6((z + 32.5) * 2)
}

// rvp6Raw converts the raw value to radar video processor values (rvp-6) by applying the
// products precision field.
func (c *Composite) rvp6Raw(value int) RVP6 {
	return RVP6(float64(value) * math.Pow10(c.precision))
}

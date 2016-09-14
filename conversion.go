package radolan

import (
	"math"
)

// Radar reflectivity factor Z in logarithmic representation dBZ: dBZ = 10 * log(Z)
type DBZ float64

// Raw radar video processor value.
type RVP6 float64

// PrecipitationRate returns the estimated precipitation rate in mm/h for the given
// reflectivity factor. The used Z-R relation is described in [6].
func (z DBZ) PrecipitationRate() float64 {
	return math.Pow(math.Pow(10, float64(z)/10)/256, 1/1.42)
}

// Reflectivity returns the estimated reflectivity factor for the given precipitation
// rate (mm/h). The used Z-R relation is described in [6].
func Reflectivity(rate float64) DBZ {
	return DBZ(10 * math.Log10(256*math.Pow(rate, 1.42)))
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

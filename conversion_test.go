package radolan

import (
	"math"
	"testing"
)

func TestConversion(t *testing.T) {
	testcases := []struct {
		rvp float32
		dbz float32
		zr  float64
	}{
		{0, -32.5, 0.0001},
		{65, 0, 0.0201},
		{100, 17.5, 0.3439},
		{200, 67.5, 1141.7670},
	}

	for _, test := range testcases {
		dbz := toDBZ(test.rvp)
		zr := PrecipitationRate(Aniol80, dbz)
		rz := Reflectivity(Aniol80, zr)
		rvp := toRVP6(dbz)

		if dbz != test.dbz {
			t.Errorf("toDBZ(%f) = %f; expected: %f", test.rvp, dbz, test.dbz)
		}
		if rvp != test.rvp {
			t.Errorf("toRVP6(toDBZ(%f)) = %f; expected: %f", test.rvp, rvp, test.rvp)
		}
		if math.Abs(test.zr-zr) > 0.0001 {
			t.Errorf("PrecipitationRate(Aniol80, toDBZ(%f)) = %f; expected: %f", test.rvp, zr, test.zr)
		}
		if math.Abs(float64(test.dbz-rz)) > 0.0000001 {
			t.Errorf("Reflectivity(PrecipitationRate(toDBZ(%f))) = %f; expected: %f",
				test.rvp, rz, test.dbz)
		}
	}
}

package radolan

import (
	"math"
	"testing"
)

func TestConversion(t *testing.T) {
	testcases := []struct {
		rvp  RVP6
		dbz  DBZ
		rate float64
	}{
		{0, -32.5, 0.0001},
		{65, 0, 0.0201},
		{100, 17.5, 0.3439},
		{200, 67.5, 1141.7670},
	}

	for _, test := range testcases {
		dbz := test.rvp.ToDBZ()
		rate := dbz.PrecipitationRate()
		rvp := dbz.ToRVP6()

		if dbz != test.dbz {
			t.Errorf("RVP6(%f).ToDBZ() = %f; expected: %f", test.rvp, dbz, test.dbz)
		}
		if rvp != test.rvp {
			t.Errorf("RVP6(%f).ToDBZ().ToRVP6() = %f; expected: %f", test.rvp, rvp, test.rvp)
		}
		if math.Abs(test.rate-rate) > 0.0001 {
			t.Errorf("RVP6(%f).ToDBZ().PrecipitationRate() = %f; expected: %f", test.rvp, rate, test.rate)
		}
	}
}
package radolan

import (
	"math"
	"testing"
)

func TestTranslate(t *testing.T) {
	equal := func(a, b float64) bool {
		epsilon := 0.1 // inaccuracy by 100 meters
		return math.Abs(a-b) < epsilon
	}

	nationalGridPG := [][]float64{
		[]float64{54.6547, 01.9178, 0.0, 0.0},
		[]float64{54.8110, 15.8706, 1.0, 0.0},
		[]float64{51.0000, 09.0000, 0.5, 0.5},
		[]float64{46.8695, 03.4921, 0.0, 1.0},
		[]float64{46.9894, 14.7218, 1.0, 1.0},
	}

	nationalGrid := [][]float64{
		[]float64{54.5877, 02.0715, 0.0, 0.0},
		[]float64{54.7405, 15.7208, 1.0, 0.0},
		[]float64{51.0000, 09.0000, 0.5, 0.5},
		[]float64{46.9526, 03.5889, 0.0, 1.0},
		[]float64{47.0705, 14.6209, 1.0, 1.0},
	}

	extendedGrid := [][]float64{
		[]float64{56.5423, -0.8654, 0.0, 0.0},
		[]float64{56.4505, 21.6986, 1.0, 0.0},
		[]float64{51.0000, 09.0000, 3 / 7.0, 7 / 15.0},
		[]float64{43.9336, 02.3419, 0.0, 1.0},
		[]float64{43.8736, 18.2536, 1.0, 1.0},
	}

	dummyPG := NewDummy("PG", 460, 460)
	dummyFZ := NewDummy("FZ", 450, 450)
	dummyRX := NewDummy("RX", 900, 900)
	dummyEX := NewDummy("EX", 1400, 1500)

	testcases := []struct {
		comp *Composite
		edge [][]float64
	}{
		{dummyPG, nationalGridPG},
		{dummyFZ, nationalGrid},
		{dummyRX, nationalGrid},
		{dummyEX, extendedGrid},
	}

	for _, test := range testcases {
		t.Logf("dummy%s: Rx = %f; Ry = %f\n",
			test.comp.Product, test.comp.Rx, test.comp.Ry)
		t.Logf("dummy%s: offx = %f; offy = %f\n",
			test.comp.Product, test.comp.offx, test.comp.offy)

		for _, edge := range test.edge {
			// result
			rx, ry := test.comp.Translate(edge[0], edge[1])
			//expected
			ex := float64(test.comp.Dx-1) * edge[2]
			ey := float64(test.comp.Dy-1) * edge[3]

			if !equal(rx, ex) || !equal(ry, ey) {
				t.Errorf("dummy%s.Translate(%#v, %#v) = (%#v, %#v); expected: (%#v, %#v)",
					test.comp.Product, edge[0], edge[1], rx, ry, ex, ey)
			}
		}
	}
}

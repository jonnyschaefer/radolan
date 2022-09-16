package radolan

import (
	"testing"
)

func Test_DE1200_WGS84(t *testing.T) {
	DE1200Grid := [][]float64{
		[]float64{55.86208711, 1.463301510, 0, 0},       // NW
		[]float64{55.84543856, 18.73161645, 1100, 0},    // NE
		[]float64{45.68460578, 16.58086935, 1100, 1200}, // SE
		[]float64{45.69642538, 3.566994635, 0, 1200},    // SW
	}

	comp := NewDummy("WN", 5, 1100, 1200)

	for _, v := range DE1200Grid {
		rx, ry := comp.Project(v[0], v[1])
		ex, ey := v[2], v[3]

		if dist(rx, ry, ex, ey) > 0.000001 {
			t.Errorf("comp.Project(%#v, %#v) = (%#v, %#v); expected: (%#v, %#v)", v[0], v[1], rx, ry, ex, ey)
		}
	}
}

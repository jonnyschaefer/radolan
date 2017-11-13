package radolan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"testing"
)

func absequal(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func dist(x1, y1, x2, y2 float64) float64 {
	x := x1 - x2
	y := y1 - y2
	return math.Sqrt(x*x + y*y)
}

func TestResolution(t *testing.T) {
	var (
		srcLat, srcLon = 48.173146, 11.546604 // Munich
		dstLat, dstLon = 53.534366, 08.576135 // Bremerhaven
		expDist        = 663.629945199998     // km
	)

	dummys := []*Composite{
		NewDummy("SF", 900, 900),
		NewDummy("SF", 450, 450),
		NewDummy("SF", 225, 225),
		NewDummy("SF", 112, 112),
		NewDummy("WX", 900, 1100),
		NewDummy("WX", 450, 550),
		NewDummy("EX", 1400, 1500),
	}

	for _, comp := range dummys {
		srcX, srcY := comp.Translate(srcLat, srcLon)
		dstX, dstY := comp.Translate(dstLat, dstLon)

		resDist := dist(srcX*comp.Rx, srcY*comp.Ry, dstX*comp.Rx, dstY*comp.Ry)

		if !absequal(resDist, expDist, 0.000001) { // inaccuracy by 1mm
			t.Errorf("dummy.Rx = %#v, dummy.Ry = %#v; distance: %#v expected: %#v)",
				comp.Rx, comp.Ry, resDist, expDist)
		}
	}

}

func TestTranslate(t *testing.T) {
	nationalGridPG := [][]float64{
		[]float64{54.6547, 01.9178, 0, 0},
		[]float64{54.8110, 15.8706, 460, 0},
		[]float64{51.0000, 09.0000, 230, 230},
		[]float64{46.8695, 03.4921, 0, 460},
		[]float64{46.9894, 14.7218, 460, 460},
	}

	nationalGridHalf := [][]float64{
		[]float64{54.5877, 02.0715, 0, 0},
		[]float64{54.7405, 15.7208, 450, 0},
		[]float64{51.0000, 09.0000, 225, 225},
		[]float64{46.9526, 03.5889, 0, 450},
		[]float64{47.0705, 14.6209, 450, 450},
	}

	nationalGrid := [][]float64{
		[]float64{54.5877, 02.0715, 0, 0},
		[]float64{54.7405, 15.7208, 900, 0},
		[]float64{51.0000, 09.0000, 450, 450},
		[]float64{46.9526, 03.5889, 0, 900},
		[]float64{47.0705, 14.6209, 900, 900},
	}

	extendedNationalGrid := [][]float64{
		[]float64{55.5482, 03.0889, 0, 0},
		[]float64{55.5342, 17.1128, 900, 0},
		[]float64{51.0000, 09.0000, 370, 550},
		[]float64{46.1929, 04.6759, 0, 1100},
		[]float64{46.1827, 15.4801, 900, 1100},
	}

	middleEuropeanGrid := [][]float64{
		[]float64{56.5423, -0.8654, 0, 0},
		[]float64{56.4505, 21.6986, 1400, 0},
		[]float64{51.0000, 09.0000, 600, 700},
		[]float64{43.9336, 02.3419, 0, 1500},
		[]float64{43.8736, 18.2536, 1400, 1500},
	}

	dummyPG := NewDummy("PG", 460, 460)
	dummyFZ := NewDummy("FZ", 450, 450)
	dummyRX := NewDummy("RX", 900, 900)
	dummyWX := NewDummy("WX", 900, 1100)
	dummyEX := NewDummy("EX", 1400, 1500)

	testcases := []struct {
		comp *Composite
		edge [][]float64
	}{
		{dummyPG, nationalGridPG},
		{dummyFZ, nationalGridHalf},
		{dummyRX, nationalGrid},
		{dummyWX, extendedNationalGrid},
		{dummyEX, middleEuropeanGrid},
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
			ex := edge[2]
			ey := edge[3]

			// allowed inaccuracy by 100 meters
			if dist(rx, ry, ex, ey) > 0.1 {
				t.Errorf("dummy%s.Translate(%#v, %#v) = (%#v, %#v); expected: (%#v, %#v)",
					test.comp.Product, edge[0], edge[1], rx, ry, ex, ey)
			}
		}
	}
}

type gridMode string

const (
	gridBottom gridMode = "bottom"
	gridCenter          = "center"
)

func TestGrid(t *testing.T) {
	testGrid(t, gridCenter,
		NewDummy("SF", 900, 900),
		NewDummy("EX", 1400, 1500),
		NewDummy("WX", 900, 1100),
	)

	testGrid(t, gridBottom,
		NewDummy("SF", 900, 900),
		NewDummy("EX", 1400, 1500),
		// NewDummy("WX", 900, 1100), testdata unavailable
	)
}

func testGrid(t *testing.T, mode gridMode, dummys ...*Composite) {
	t.Helper()

	var offx, offy float64
	switch mode {
	case gridBottom:
		offx, offy = 0.0, 1.0
	case gridCenter:
		offx, offy = 0.5, 0.5
	default:
		t.Fatalf("unknown grid mode %#v", mode)
	}

	for _, comp := range dummys {
		lname := fmt.Sprintf("testdata/lambda_%s_%dx%d.txt", mode, comp.Dy, comp.Dx)
		pname := fmt.Sprintf("testdata/phi_%s_%dx%d.txt", mode, comp.Dy, comp.Dx)

		lbuf, err := ioutil.ReadFile(lname)
		if err != nil {
			t.Fatal(err)
		}
		pbuf, err := ioutil.ReadFile(pname)
		if err != nil {
			t.Fatal(err)
		}

		// fortran format F8.5 means read 8 bytes, the last 5 bytes are decimal places
		const length = 8

		var l, p int64

		// Beschreibung-E-Produkte-Raster.pdf Radolan-Cons Version 1.0:
		// "Die Dateien beginnen mit dem Referenzwert des Datenelements in der linken unteren Ecke des
		//  Komposits, spaltenweise von Westen nach Osten und zeilenweise von Süden nach Norden."
		//
		// "Für die Dateien mit der Bezeichnung _bottom beziehen sich die Koordinaten
		//  jeweils auf die linke untere Ecke jedes Datenelements. Für Dateien mit der Bezeichnung
		//  _center auf den Zentralpunkt."

		for y := comp.Dy - 1; y >= 0; y-- {
			for x := 0; x < comp.Dx; x++ {

				// newlines can occur in input files.
				for lbuf[l] == '\n' {
					l++
				}
				for pbuf[p] == '\n' {
					p++
				}

				lstring := string(bytes.TrimSpace(lbuf[l : l+length]))
				lamda, err := strconv.ParseFloat(lstring, 64)
				if err != nil {
					t.Fatalf("invalid grid coordinate at (%d, %d): %#v %s", x, y, lstring, lname)
				}
				l += length

				pstring := string(bytes.TrimSpace(pbuf[p : p+length]))
				phi, err := strconv.ParseFloat(pstring, 64)
				if err != nil {
					t.Fatalf("invalid grid coordinate at (%d, %d): %#v %s", x, y, pstring, pname)
				}
				p += length

				tx, ty := comp.Translate(phi, lamda)
				ex, ey := float64(x)+offx, float64(y)+offy
				if dist(tx, ty, ex, ey) > 0.01 { // 10m
					t.Errorf("dummy%s.Translate(%#v, %#v) = (%#v, %#v); expected: (%#v, %#v)",
						comp.Product, phi, lamda, tx, ty, ex, ey)
				}
			}
		}

		if len(bytes.TrimSpace(lbuf[l:])) > 0 {
			t.Fatalf("unprocessed data remaining in %s", lname)
		}
		if len(bytes.TrimSpace(pbuf[p:])) > 0 {
			t.Fatalf("unprocessed data remaining in %s", pname)
		}
	}
}

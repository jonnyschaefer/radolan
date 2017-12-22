package radolan

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

type headerTestcase struct {
	// head of file
	test string

	// expected
	expBinary       string
	expProduct      string
	expCaptureTime  time.Time
	expForecastTime time.Time
	expInterval     time.Duration
	expDx           int
	expDy           int
	expDataLength   int
	expPrecision    int
	expLevel        []float32
}

func TestParseHeaderPG(t *testing.T) {
	ht := &headerTestcase{}
	var err1, err2 error

	// head of file
	ht.test = "PG262115100000616BY22205LV 6  1.0 19.0 28.0 37.0 46.0 55.0CS0MX 0MS " +
		"88<boo,ros,emd,hnr,umd,pro,ess,fld,drs,neu,nhb,oft,eis,tur,isn,fbg,mem> " +
		"are used, BG460460\x03binarycontent"

	// expected
	ht.expBinary = "binarycontent"
	ht.expProduct = "PG"
	ht.expCaptureTime, err1 = time.Parse(time.RFC1123, "Sun, 26 Jun 2016 23:15:00 CEST")
	ht.expForecastTime, err2 = time.Parse(time.RFC1123, "Sun, 26 Jun 2016 23:15:00 CEST")
	ht.expDx = 460
	ht.expDy = 460
	ht.expDataLength = 22205 - 159 // BY - header_etx_length
	ht.expPrecision = 0
	ht.expLevel = []float32{1.0, 19.0, 28.0, 37.0, 46.0, 55.0}

	if err1 != nil || err2 != nil {
		t.Errorf("%s.parseHeader(): wrong testcase time.Parse", ht.expProduct)
	}

	testParseHeader(t, ht)
}

func TestParseHeaderFZ(t *testing.T) {
	ht := &headerTestcase{}
	var err1, err2 error

	// head of file
	ht.test = "FZ282105100000716BY 405160VS 3SW   2.13.1PR E-01INT   5GP 450x 450VV 100MF " +
		"00000002MS 66<boo,ros,emd,hnr,umd,pro,ess,drs,neu,nhb,oft,eis,tur,isn,fbg,mem>" +
		"\x03binarycontent"

	// ht.expected values
	ht.expBinary = "binarycontent"

	ht.expProduct = "FZ"
	ht.expCaptureTime, err1 = time.Parse(time.RFC1123, "Thu, 28 Jul 2016 23:05:00 CEST")
	ht.expForecastTime, err2 = time.Parse(time.RFC1123, "Fri, 29 Jul 2016 00:45:00 CEST")
	ht.expInterval = 5 * time.Minute
	ht.expDx = 450
	ht.expDy = 450
	ht.expDataLength = 405160 - 154 // BY - header_etx_length
	ht.expPrecision = -1
	ht.expLevel = []float32(nil)

	if err1 != nil || err2 != nil {
		t.Errorf("%s.parseHeader(): wrong testcase time.Parse", ht.expProduct)
	}

	testParseHeader(t, ht)
}

func testParseHeader(t *testing.T, ht *headerTestcase) {
	dummy := &Composite{}
	reader := bufio.NewReader(strings.NewReader(ht.test))

	// run
	if err := dummy.parseHeader(reader); err != nil {
		t.Errorf("%s.parseHeader(): returned error: %#v", err.Error())
	}

	// test results
	// Product
	if dummy.Product != ht.expProduct {
		t.Errorf("%s.parseHeader(): Product: %#v; expected: %#v", ht.expProduct,
			dummy.Product, ht.expProduct)
	}

	// CaptureTime
	if !dummy.CaptureTime.Equal(ht.expCaptureTime) {
		t.Errorf("%s.parseHeader(): CaptureTime: %#v; expected: %#v", ht.expProduct,
			dummy.CaptureTime.String(), ht.expCaptureTime.String())
	}

	// ForecastTime
	if !dummy.ForecastTime.Equal(ht.expForecastTime) {
		t.Errorf("%s.parseHeader(): ForecastTime: %#v; expected: %#v", ht.expProduct,
			dummy.ForecastTime.String(), ht.expForecastTime.String())
	}

	// Interval
	if dummy.Interval != ht.expInterval {
		t.Errorf("%s.parseHeader(): Interval: %#v; expected: %#v", ht.expProduct,
			dummy.Interval.String(), ht.expInterval.String())
	}

	// Dx Dy
	if dummy.Dx != ht.expDx || dummy.Dy != ht.expDy {
		t.Errorf("%s.parseHeader(): Dx: %d Dy: %d; expected Dx: %d Dy: %d", ht.expProduct,
			dummy.Dx, dummy.Dy, ht.expDx, ht.expDy)
	}

	// dataLength
	if dummy.dataLength != ht.expDataLength {
		t.Errorf("%s.parseHeader(): dataLength: %#v; expected: %#v", ht.expProduct,
			dummy.dataLength, ht.expDataLength)
	}

	// precision
	if dummy.precision != ht.expPrecision {
		t.Errorf("%s.parseHeader(): precision: %#v; expected: %#v", ht.expProduct,
			dummy.precision, ht.expPrecision)
	}

	// level
	for i := range ht.expLevel {
		if len(dummy.level) != len(ht.expLevel) || dummy.level[i] != ht.expLevel[i] {
			t.Errorf("%s.parseHeader(): level: %#v; expected: %#v", ht.expProduct,
				dummy.level, ht.expLevel)
		}
	}

	// check consistency
	if line, _ := reader.ReadString('\n'); line != ht.expBinary {
		t.Errorf("%s.parseHeader(): binary data corrupted", ht.expProduct)
	}
}

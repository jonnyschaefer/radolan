package radolan

import (
	"bufio"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// splitHeader splits the given header string into its fields. The returned
// map is using the field name as key and the field content as value.
func splitHeader(header string) (m map[string]string) {
	m = make(map[string]string)
	var beginKey, endKey, beginValue, endValue int
	var dispatch bool

	for i, c := range header {
		if unicode.IsUpper(c) {
			if dispatch {
				m[header[beginKey:endKey]] = header[beginValue:endValue]
				beginKey = i
				dispatch = false
			}
			endKey = i + 1
		} else {
			if i == 0 {
				return // no key prefixing value
			}
			if !dispatch {
				beginValue = i
				dispatch = true
			}
			endValue = i + 1
		}
	}
	m[header[beginKey:endKey]] = header[beginValue:endValue]

	return
}

// parseHeader parses and the composite header and writes the related fields as
// described in [1] and [3].
func (c *Composite) parseHeader(reader *bufio.Reader) error {
	header, err := reader.ReadString('\x03')
	if err != nil || len(header) < 22 { // smaller length makes no sense
		return newError("parseHeader", "header corrupted: too short")
	}

	// Split header segments
	section := splitHeader(header[:len(header)-1]) // without delimiter

	// Parse Product - Example: "PG" or "FZ"
	c.Product = header[:2]

	// Parse DataLength - Example: "BY 405160"
	if _, err := fmt.Sscanf(section["BY"], "%d", &c.dataLength); err != nil {
		return newError("parseHeader", "could not parse data length: "+err.Error())
	}
	c.dataLength -= len(header) // remove header length including delimiter

	// Parse CaptureTime - Example: "PG262115100000616" or "FZ211615100000716"
	date := header[2:8] + header[13:17] // cut WMO number
	c.CaptureTime, err = time.Parse("0215040106", date)
	if err != nil {
		return newError("parseHeader", "could not parse capture time: "+err.Error())
	}

	// Parse ForecastTime - Example: "VV 005"
	c.ForecastTime = c.CaptureTime
	if vv, ok := section["VV"]; ok {
		min := 0
		if _, err := fmt.Sscanf(vv, "%d", &min); err != nil {
			return newError("parseHeader", "could not parse forecast time: "+err.Error())
		}
		c.ForecastTime = c.CaptureTime.Add(time.Duration(min) * time.Minute)
	}

	// Parse Dimensions - Example: "GP 450x 450" or "BG460460" or "GP 1500x1400"
	dim := section["GP"]
	if bg, ok := section["BG"]; ok {
		dim = bg[:len(bg)/2] + "x" + bg[len(bg)/2:]
	}

	if _, err := fmt.Sscanf(dim, "%dx%d", &c.Dy, &c.Dx); err != nil {
		return newError("parseHeader", "could not parse dimensions: "+err.Error())
	}

	// Parse Precision - Example: "PR E-01" or "PR E+00"
	if prec, ok := section["E"]; ok { // not that nice
		if _, err := fmt.Sscanf(prec, "%d", &c.precision); err != nil {
			return newError("parseHeader", "could not parse precision: "+err.Error())
		}
	}

	// Parse Level - Example "LV 6  1.0 19.0 28.0 37.0 46.0 55.0"
	if lv, ok := section["LV"]; ok {
		level := strings.Fields(lv)
		if len(level) < 2 {
			return newError("parseHeader", "invalid level count")
		}

		c.level = make([]DBZ, len(level)-1)
		for i, f := range level[1:] {
			if _, err = fmt.Sscanf(f, "%f", &c.level[i]); err != nil {
				return newError("parseHeader", "invalid level value: "+err.Error())
			}
		}
	}

	return nil
}

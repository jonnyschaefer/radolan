package radolan

import (
	"bufio"
	"math"
)

// parseRunlength parses the runlength encoded composite and writes into the
// previously created PlainData field of the composite.
func (c *Composite) parseRunlength(reader *bufio.Reader) error {
	for i := range c.PlainData {
		line, err := c.readLineRunlength(reader)
		if err != nil {
			return err
		}

		err = c.decodeRunlength(c.PlainData[i], line)
		if err != nil {
			return err
		}
	}

	return nil
}

// readLineRunlength reads a line until newline (non inclusive) from the given reader.
// This method is used to get a line of runlenth encoded data.
func (c *Composite) readLineRunlength(rd *bufio.Reader) (line []byte, err error) {
	line, err = rd.ReadBytes('\x0A')
	if err != nil {
		err = newError("readLineRunlength", err.Error())
	}
	length := len(line)
	if length > 0 {
		line = line[:length-1]
	}
	return
}

// decodeRunlength decodes the source line and writes to the given destination.
func (c *Composite) decodeRunlength(dst []RVP6, line []byte) error {
	// fill destination as runlength encoding will induce gaps
	nan := RVP6(math.NaN())
	for i := range dst {
		dst[i] = nan
	}

	dstpos := 0
	offset := true
	for i, value := range line {
		switch true {
		case i == 0: // skip useless line number
		case offset: // calculate offset
			if value < 16 {
				return newError("decodeRunlength", "invalid offset value")
			}

			dstpos += int(value) - 16 // update offset position
			offset = value == 255     // see if next byte will be also offset
		default:
			// value [XXXX|YYYY] decodes to YYYY repeated XXXX times.
			runlength := int(value >> 4)
			value &= 0x0F

			for j := 0; j < runlength; j++ {
				if dstpos >= len(dst) {
					return newError("decodeRunlength", "destination size exceeded")
				}

				dst[dstpos] = c.rvp6Runlength(value)
				dstpos++
			}
		}
	}

	return nil
}

// rvp6Runlength sets the value of level based composite products to radar
// video processor values (rvp-6).
func (c *Composite) rvp6Runlength(value byte) RVP6 {
	if value == 0 {
		return RVP6(math.NaN())
	}
	value--

	if int(value) >= len(c.level) { // border markings
		return RVP6(math.NaN())
	}
	return c.level[value]
}

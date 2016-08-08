package radolan

import (
	"bufio"
	"io"
	"math"
)

// parseSingleByte parses the single byte encoded composite as described in [1] and writes
// into the previously created Data field of the composite.
func (c *Composite) parseSingleByte(reader *bufio.Reader) error {
	last := len(c.Data) - 1
	for i := range c.Data {
		line, err := c.readLineSingleByte(reader)
		if err != nil {
			return err
		}

		err = c.decodeSingleByte(c.Data[last-i], line) // write vertically flipped
		if err != nil {
			return err
		}
	}

	return nil
}

// readLineSingleByte reads a line until horizontal limit from the given reader
// This method is used to get a line of single byte encoded data.
func (c *Composite) readLineSingleByte(rd *bufio.Reader) (line []byte, err error) {
	line = make([]byte, c.Dx)
	_, err = io.ReadFull(rd, line)
	if err != nil {
		err = newError("readLineSingleByte", err.Error())
	}
	return
}

// decodeSingleByte decodes the source line and writes to the given destination.
func (c *Composite) decodeSingleByte(dst []RVP6, line []byte) error {
	if len(dst) != len(line) {
		return newError("decodeSingleByte", "wrong destination or source size")
	}

	for i, v := range line {
		dst[i] = c.rvp6SingleByte(v)
	}

	return nil
}

// rvp6SingleByte converts the raw byte of single byte encoded
// composite products to radar video processor values (rvp-6). NaN may be returned
// when the no-data flag is set.
func (c *Composite) rvp6SingleByte(value byte) RVP6 {
	if value == 250 { // error code: no-data
		return RVP6(math.NaN())
	}

	// set decimal point
	return c.rvp6Raw(int(value))
}

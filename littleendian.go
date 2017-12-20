package radolan

import (
	"bufio"
	"io"
	"math"
)

// parseLittleEndian parses the little endian encoded composite as described in [1] and [3].
// Result are written into the previously created PlainData field of the composite.
func (c *Composite) parseLittleEndian(reader *bufio.Reader) error {
	last := len(c.PlainData) - 1
	for i := range c.PlainData {
		line, err := c.readLineLittleEndian(reader)
		if err != nil {
			return err
		}

		err = c.decodeLittleEndian(c.PlainData[last-i], line) // write vertically flipped
		if err != nil {
			return err
		}
	}

	return nil
}

// readLineLittleEndian reads a line until horizontal limit from the given reader
// This method is used to get a line of little endian encoded data.
func (c *Composite) readLineLittleEndian(rd *bufio.Reader) (line []byte, err error) {
	line = make([]byte, c.Dx*2)
	_, err = io.ReadFull(rd, line)
	if err != nil {
		err = newError("readLineLittleEndian", err.Error())
	}
	return
}

// decodeLittleEndian decodes the source line and writes to the given destination.
func (c *Composite) decodeLittleEndian(dst []RVP6, line []byte) error {
	if len(line)%2 != 0 || len(dst)*2 != len(line) {
		return newError("decodeLittleEndian", "wrong destination or source size")
	}

	for i := range dst {
		tuple := [2]byte{line[2*i], line[2*i+1]}
		dst[i] = c.rvp6LittleEndian(tuple)
	}

	return nil
}

// rvp6LittleEndian converts the raw two byte tuple of little endian encoded composite products
// to radar video processor values (rvp-6). NaN may be returned when the no-data flag is set.
func (c *Composite) rvp6LittleEndian(tuple [2]byte) RVP6 {
	var value int = 0x0F & int(tuple[1])
	value = (value << 8) + int(tuple[0])

	if tuple[1]&(1<<5) != 0 { // error code: no-data
		return RVP6(math.NaN())
	}

	if tuple[1]&(1<<6) != 0 { // flag: negative value
		value *= -1
	}

	// set decimal point
	return c.rvp6Raw(value)
}

package radolan

import (
	"bufio"
)

// encoding types of the composite
type encoding int

const (
	runlength encoding = iota
	littleEndian
	singleByte
	unknown
)

// parsing methods
var parse = [4]func(c *Composite, rd *bufio.Reader) error{}

// init maps the parsing methods to the encoding type
func init() {
	parse[runlength] = (*Composite).parseRunlength
	parse[littleEndian] = (*Composite).parseLittleEndian
	parse[singleByte] = (*Composite).parseSingleByte
	parse[unknown] = (*Composite).parseUnknown
}

// identifyEncoing identifies the encoding type of the data section by
// only comparing header characteristics.
// This method requires header data to be already written.
func (c *Composite) identifyEncoding() encoding {
	values := c.Dx * c.Dy

	if c.level != nil {
		return runlength
	}
	if c.dataLength == values*2 {
		return littleEndian
	}
	if c.dataLength == values {
		return singleByte
	}

	return unknown
}

// parseData parses the composite data and writes the related fields.
// This method requires header data to be already written.
func (c *Composite) parseData(reader *bufio.Reader) error {
	if c.Dx == 0 || c.Dy == 0 {
		return newError("parseData", "parsed header data required")
	}

	// create Data fields
	c.Data = make([][]RVP6, c.Dy)
	for i := range c.Data {
		c.Data[i] = make([]RVP6, c.Dx)
	}

	return parse[c.identifyEncoding()](c, reader)
}

// parseUnknown performs no action and always returns an error.
func (c *Composite) parseUnknown(rd *bufio.Reader) error {
	return newError("parseUnknown", "unknown encoding")
}

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
	values := c.Px * c.Py

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
	if c.Px == 0 || c.Py == 0 {
		return newError("parseData", "parsed header data required")
	}

	// create Data fields
	c.PlainData = make([][]RVP6, c.Py)
	for i := range c.PlainData {
		c.PlainData[i] = make([]RVP6, c.Px)
	}

	return parse[c.identifyEncoding()](c, reader)
}

// arrangeData slices plain data into its data layers or strips preceeding
// vertical projection
func (c *Composite) arrangeData() {
	if c.Py%c.Dy == 0 { // multiple layers are linked downwards
		c.DataZ = make([][][]RVP6, c.Py/c.Dy)
		for i := range c.DataZ {
			c.DataZ[i] = c.PlainData[c.Dy*i : c.Dy*(i+1)] // split layers
		}
	} else { // only use bottom most part of plain data
		c.DataZ = [][][]RVP6{c.PlainData[c.Py-c.Dy:]} // strip elevation
	}

	c.Dz = len(c.DataZ)
	c.Data = c.DataZ[0] // alias
}

// parseUnknown performs no action and always returns an error.
func (c *Composite) parseUnknown(rd *bufio.Reader) error {
	return newError("parseUnknown", "unknown encoding")
}

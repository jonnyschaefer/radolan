// radolan2png is an example program for the radolan package, that converts
// radolan composite files to .png images. The created images also contain
// an overlay showing the german borders and a latitude longitude mesh.
package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"gitlab.cs.fau.de/since/radolan"
)

var (
	borderColor = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
	meshColor   = color.RGBA{0x33, 0xFF, 0x22, 0xFF}
)

func main() {
	// display help message
	if len(os.Args) < 3 {
		fmt.Printf("radolan2png converts radolan composite files to png images."+
			"\n\n\tUsage: %s <input> <output.png>\n\n", os.Args[0])
		return
	}

	convert(os.Args[1], os.Args[2])
}

func convert(in, out string) {
	// open input file
	infile, err := os.Open(in)
	care(err)
	defer infile.Close()

	// create new composite
	comp, err := radolan.NewComposite(infile)
	care(err)

	// choose color function
	heatmap := radolan.HeatmapReflectivity
	switch comp.Product[0] {
	case 'R', 'W', 'S', 'E':
		if comp.Product[1] != 'X' {
			heatmap = radolan.HeatmapAccumulatedDay
		}
	}

	// convert composite to image using the color function
	img := comp.Image(heatmap)

	// draw borders
	for _, b := range border {
		// convert border points to data indices
		x, y := comp.Translate(b[0], b[1])

		// draw point
		img.Set(int(x), int(y), borderColor)
	}

	// draw mesh
	for e := 1.0; e < 16.0; e += 0.1 {
		for n := 46.0; n < 55.0; n += 0.1 {
			if e-float64(int(e)) < 0.1 || n-float64(int(n)) < 0.1 {
				x, y := comp.Translate(n, e)
				img.Set(int(x), int(y), meshColor)
			}
		}
	}

	// create output file
	outfile, err := os.Create(out)
	care(err)
	defer outfile.Close()

	// write image to output file
	care(png.Encode(outfile, img))
}

// care exits the program if an error occured
func care(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

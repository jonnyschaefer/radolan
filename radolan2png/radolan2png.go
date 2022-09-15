// radolan2png is an example program for the radolan package, that converts
// radolan composite files to .png images. The created images also contain
// an overlay showing the german borders and a latitude longitude mesh.
package main

import (
	"fmt"
	"gitlab.cs.fau.de/since/radolan"
	"gitlab.cs.fau.de/since/radolan/radolan2png/vis"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
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

	fmt.Printf("%s-Image (%s) showing %s\n", comp.Product, comp.DataUnit, comp.ForecastTime)

	var heatmap vis.ColorFunc
	switch comp.DataUnit {
	case radolan.Unit_mm:
		max := 200.0
		if comp.Interval <= time.Hour {
			max = 100.0
		}
		if comp.Interval >= time.Hour*24*7 {
			max = 400.0
		}
		heatmap = vis.Heatmap(0.1, max, vis.Log)
	case radolan.Unit_dBZ:
		heatmap = vis.HeatmapReflectivity
	case radolan.Unit_km:
		heatmap = vis.Graymap(0, 15, vis.Id)
	case radolan.Unit_mps:
		heatmap = vis.HeatmapRadialVelocity
	}

	// convert composite to image using the color function
	img := vis.Image(heatmap, comp, 0) // TODO: select layer

	// draw borders
	if comp.HasProjection {
		// print grid dimensions
		fmt.Printf("detected grid: %.1f km * %.1f km\n", float64(comp.Dx)*comp.Rx, float64(comp.Dy)*comp.Ry)

		for _, b := range border {
			// convert border points to data indices
			x, y := comp.Project(b[0], b[1])

			// draw point
			img.Set(int(x), int(y), borderColor)
		}

		// draw mesh
		for e := 1.0; e < 16.0; e += 0.1 {
			for n := 46.0; n < 55.0; n += 0.1 {
				if e-float64(int(e)) < 0.1 || n-float64(int(n)) < 0.1 {
					x, y := comp.Project(n, e)
					img.Set(int(x), int(y), meshColor)
				}
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

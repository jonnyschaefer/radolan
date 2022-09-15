package vis

import (
	"gitlab.cs.fau.de/since/radolan"
	"image"
	"image/color"
	"math"
)

// A ColorFunc can be used to assign colors to data values for image creation.
type ColorFunc func(val float64) color.RGBA

// Sample color and grayscale gradients for visualization with the image method.
var (
	// HeatmapReflectivityShort is a color gradient for cloud reflectivity
	// composites between 5dBZ and 75dBZ.
	HeatmapReflectivityShort = Heatmap(5.0, 75.0, Id)

	// HeatmapReflectivity is a color gradient for cloud reflectivity
	// composites between 5dBZ and 75dBZ.
	HeatmapReflectivity = Heatmap(1.0, 75.0, Id)

	// HeatmapReflectivityWide is a color gradient for cloud reflectivity
	// composites between -32.5dBZ and 75dBZ.
	HeatmapReflectivityWide = Heatmap(-32.5, 75.0, Id)

	// HeatmapAccumulatedHour is a color gradient for accumulated rainfall
	// composites (e.g RW) between 0.1mm/h and 100 mm/h using logarithmic
	// compression.
	HeatmapAccumulatedHour = Heatmap(0.1, 100, Log)

	// HeatmapAccumulatedDay is a color gradient for accumulated rainfall
	// composites (e.g. SF) between 0.1mm and 200mm using logarithmic
	// compression.
	HeatmapAccumulatedDay = Heatmap(0.1, 200, Log)

	HeatmapRadialVelocity = Radialmap(-31.5, 31.5, Log)

	// GraymapLinear is a linear grayscale gradient between the (raw) rvp-6
	// values 0 and 409.5.
	GraymapLinear = Graymap(0, 409.5, Id)

	// GraymapLinearWide is a linear grayscale gradient between the (raw)
	// rvp-6 values 0 and 4095.
	GraymapLinearWide = Graymap(0, 4095, Id)
)

// Id is the identity (no compression)
func Id(x float64) float64 {
	return x
}

// Log is the natural logarithm (logarithmic compression)
func Log(x float64) float64 {
	return math.Log(x)
}

// Image creates an image by evaluating the color function fn for each data
// value in the given z-layer.
func Image(fn ColorFunc, c *radolan.Composite, layer int) *image.RGBA {
	rec := image.Rect(0, 0, c.Dx, c.Dy)
	img := image.NewRGBA(rec)

	if layer < 0 || layer >= c.Dz {
		return img
	}

	for y := 0; y < c.Dy; y++ {
		for x := 0; x < c.Dx; x++ {
			img.Set(x, y, fn(float64(c.DataZ[layer][y][x])))
		}
	}

	return img
}

// Graymap returns a grayscale gradient between min and max. A compression function is used to
// make logarithmic scales possible.
func Graymap(min, max float64, compression func(float64) float64) ColorFunc {
	min = compression(min)
	max = compression(max)

	return func(val float64) color.RGBA {
		val = compression(val)

		if val < min {
			return color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
		}

		p := (val - min) / (max - min)
		if p > 1 {
			p = 1
		}

		l := uint8(0xFF * p)
		return color.RGBA{l, l, l, 0xFF}
	}
}

// Radialmap returns a dichromatic gradient from min to 0 to max which can
// be used for doppler radar radial velocity products.
func Radialmap(min, max float64, compression func(float64) float64) ColorFunc {
	return func(val float64) color.RGBA {
		if val != val {
			return color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
		}

		base := math.Max(math.Abs(min), math.Abs(max))
		p := compression(math.Abs(val)) / compression(base)

		if p > 1 {
			p = 1
		}
		lev := uint8(0xFF * p)

		var non byte = 0x00
		if math.Abs(val) <= 1 {
			lev = 0xFF
			non = 0xCC
		}

		if val < 0 {
			return color.RGBA{non, lev, lev, 0xFF}
		}

		return color.RGBA{lev, non, non, 0xFF}
	}
}

// Heatmap returns a colour gradient between min and max. A compression function is used to
// make logarithmic scales possible.
func Heatmap(min, max float64, compression func(float64) float64) ColorFunc {
	min = compression(min)
	max = compression(max)

	return func(val float64) color.RGBA {
		val = compression(val)
		if val < min {
			return color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
		}

		p := (val - min) / (max - min)
		if p > 1 { // limit
			p = 1
		}
		h := math.Mod(360-(330*p)+240, 360)

		s := 1.0          // saturation
		l := 0.5*p + 0.25 // lightness

		// adapted from https://en.wikipedia.org/wiki/HSL_and_HSV#From_HSL
		c := (1 - math.Abs(2*l-1)) * s // calculate chroma

		hh := h / 60
		x := c * (1 - math.Abs(math.Mod(hh, 2)-1))

		if math.IsNaN(hh) {
			hh = -1
		}

		var rr, gg, bb float64
		switch int(hh) {
		case 0:
			rr, gg, bb = c, x, 0
		case 1:
			rr, gg, bb = x, c, 0
		case 2:
			rr, gg, bb = 0, c, x
		case 3:
			rr, gg, bb = 0, x, c
		case 4:
			rr, gg, bb = x, 0, c
		case 5:
			rr, gg, bb = c, 0, x
		}

		m := l - c/2
		r, g, b := uint8(0xFF*(rr+m)), uint8(0xFF*(gg+m)), uint8(0xFF*(bb+m))

		return color.RGBA{r, g, b, 0xFF}
	}
}

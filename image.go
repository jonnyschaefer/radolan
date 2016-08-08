package radolan

import (
	"image"
	"image/color"
	"math"
)

// A ColorFunc can be used to assign colors to rvp-6 values for image creation.
type ColorFunc func(rvp RVP6) color.RGBA

// Sample color and grayscale gradients for visualization with the image method.
var (
	// HeatmapReflectivity is a color gradient for cloud reflectivity composites
	// between 5dBZ and 75dBZ.
	HeatmapReflectivity ColorFunc

	// HeatmapReflectivityWide is a color gradient for cloud reflectivity composites
	// between -32.5dBZ and 75dBZ.
	HeatmapReflectivityWide ColorFunc

	// HeatmapAccumulatedHour is a color gradient for accumulated rainfall composites
	// (e.g RW) between 0.1mm/h and 100 mm/h using logarithmic compression.
	HeatmapAccumulatedHour ColorFunc

	// HeatmapAccumulatedDay is a color gradient for accumulated rainfall composites
	// (e.g. SF) between 0.1mm and 200mm using logarithmic compression.
	HeatmapAccumulatedDay ColorFunc

	// GraymapLinear is a linear grayscale gradient between the (raw) rvp-6
	// values 0 and 409.5.
	GraymapLinear ColorFunc

	// GraymapLinearWide is a linear grayscale gradient between the (raw) rvp-6
	// values 0 and 4095.
	GraymapLinearWide ColorFunc
)

func init() {
	// return identity (no compression)
	id := func(x RVP6) RVP6 {
		return x
	}
	// return logarithm
	log := func(x RVP6) RVP6 {
		return RVP6(math.Log(float64(x)))
	}

	HeatmapReflectivity = heatmap(DBZ(5.0).ToRVP6(), DBZ(75.0).ToRVP6(), id)
	HeatmapReflectivityWide = heatmap(DBZ(-32.5).ToRVP6(), DBZ(75.0).ToRVP6(), id)
	HeatmapAccumulatedHour = heatmap(0.1, 100, log)
	HeatmapAccumulatedDay = heatmap(0.1, 200, log)

	GraymapLinear = graymap(0, 409.5, id)
	GraymapLinearWide = graymap(0, 4095, id)
}

// Image creates an image by evaluating the color function fn for each raw rvp-6 value.
func (c *Composite) Image(fn ColorFunc) *image.RGBA {
	rec := image.Rect(0, 0, c.Dx, c.Dy)
	img := image.NewRGBA(rec)

	for y := 0; y < c.Dy; y++ {
		for x := 0; x < c.Dx; x++ {
			img.Set(x, y, fn(c.Data[y][x]))
		}
	}

	return img
}

// graymap returns a grayscale gradient between min and max. A compression function is used to
// make logarithmic scales possible.
func graymap(min, max RVP6, compression func(RVP6) RVP6) ColorFunc {
	min = compression(min)
	max = compression(max)

	return func(rvp RVP6) color.RGBA {
		rvp = compression(rvp)

		if rvp < min {
			return color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
		}

		p := float64((rvp - min) / (max - min))
		if p > 1 {
			p = 1
		}

		l := uint8(0xFF * p)
		return color.RGBA{l, l, l, 0xFF}
	}
}

// heatmap returns a colour gradient between min and max. A compression function is used to
// make logarithmic scales possible.
func heatmap(min, max RVP6, compression func(RVP6) RVP6) ColorFunc {
	min = compression(min)
	max = compression(max)

	return func(rvp RVP6) color.RGBA {
		rvp = compression(rvp)
		if rvp < min {
			return color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
		}

		p := float64((rvp - min) / (max - min))
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

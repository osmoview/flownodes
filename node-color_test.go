package flownodes

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/osmoview/flowdata"
)

func createTestImageColor(width, height int, col color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{col}, image.Point{}, draw.Src)
	return img
}

func TestColor_Execute(t *testing.T) {
	testImage := createTestImageColor(100, 100, color.RGBA{255, 255, 0, 255})

	testFrame := flowdata.Frame{
		Image: testImage,
	}

	testObjects := []flowdata.Object{
		{BBox: [4]float64{0.0, 0.0, 0.5, 0.5}},
		{BBox: [4]float64{0.5, 0.5, 0.5, 0.5}},
	}

	col := &Color{
		frame:   testFrame,
		objects: testObjects,
	}

	_, err := col.Execute(nil)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	expectedCssColor := "#FFFF00"
	bounds := col.frame.Image.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for _, obj := range testObjects {
		x := int(obj.BBox[0] * float64(width))
		y := int(obj.BBox[1] * float64(height))
		w := int(obj.BBox[2] * float64(width))
		h := int(obj.BBox[3] * float64(height))

		var r, g, b, a uint32
		var count uint32
		for i := x; i < x+w; i++ {
			for j := y; j < y+h; j++ {
				ri, gi, bi, ai := col.frame.Image.At(i, j).RGBA()
				r += ri
				g += gi
				b += bi
				a += ai
				count++
			}
		}

		rgba := color.RGBA{
			R: uint8(r / count / 256),
			G: uint8(g / count / 256),
			B: uint8(b / count / 256),
			A: uint8(a / count / 256),
		}

		if col.rgbaToHex(rgba) != expectedCssColor {
			t.Errorf("Expected color %v, but got %v", expectedCssColor, rgba)
		}
	}
}

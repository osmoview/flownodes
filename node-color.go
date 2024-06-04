package flownodes

import (
	"errors"
	"fmt"
	"image/color"
	"log"

	"github.com/osmoview/flowdata"
)

var (
	ErrCoordinatesOutsideImage = errors.New("coordinates outside the image")
)

type Color struct {
	// frame image
	frame flowdata.Frame
	// objects Bbox in relative coordinates format [x,y,w,h]
	objects []flowdata.Object

	port Port
}

func (col *Color) Init(db DB) error {
	col.port = Port{Name: "objects", Type: PortObjects}

	return nil
}

func (col *Color) Execute(c Ctx) (out Resp, err error) {
	img := col.frame.GetImage()

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for key := range col.objects {
		// convert relative coordinates to absolute ones
		x := int(col.objects[key].BBox[0] * float64(width))
		y := int(col.objects[key].BBox[1] * float64(height))
		w := int(col.objects[key].BBox[2] * float64(width))
		h := int(col.objects[key].BBox[3] * float64(height))

		// check if the coordinates are outside the image bound
		if x < 0 || y < 0 || x+w > width || y+h > height {
			log.Println(ErrCoordinatesOutsideImage)
			continue
		}

		var r, g, b, a uint32
		var cnt uint32

		for i := x; i < x+w; i++ {
			for j := y; j < y+h; j++ {
				ri, gi, bi, ai := img.At(i, j).RGBA()

				r += ri
				g += gi
				b += bi
				a += ai

				cnt++
			}
		}

		rgba := color.RGBA{
			R: uint8(r / cnt / 256),
			G: uint8(g / cnt / 256),
			B: uint8(b / cnt / 256),
			A: uint8(a / cnt / 256),
		}

		col.objects[key].Color = col.rgbaToHex(rgba)
	}

	return Resp{col.port.Name: col.objects}, nil
}

func (col *Color) GetPorts() (in, out []Port) {
	return in, []Port{col.port}
}

// rgbaToHex convert a color rgba to a hex string
func (col *Color) rgbaToHex(c color.RGBA) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}

	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

package flownodes

import (
	"errors"
	"image"
	"image/draw"
	"log"

	"github.com/osmoview/flowdata"

	"github.com/EdlinOrg/prominentcolor"
)

type Color struct {
	// frame image
	frame Port
	// objects Bbox in relative coordinates format [x,y,w,h]
	objects Port
}

func (col *Color) Init(db DB) (err error) {
	col.frame = Port{Name: PortFrame, Type: PortFrame}
	col.objects = Port{Name: PortObjects, Type: PortObjects}

	return
}

func (col *Color) Execute(c Ctx) (out Resp, err error) {
	var frame *flowdata.Frame

	if data, ok := c.GetData(PortFrame); ok {
		frame, ok = data.(*flowdata.Frame)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
	}
	if frame == nil {
		return nil, errors.New("frame not found")
	}

	objects, ok := c.MustGetData(PortObjects).([]*flowdata.Object)
	if !ok {
		return nil, ErrUnexpectedDataType
	}

	newObjects := make([]*flowdata.Object, 0, len(objects))

	img := frame.GetImage()
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for _, obj := range objects {
		x := int(obj.BBox[0] * float64(width))
		y := int(obj.BBox[1] * float64(height))
		w := int(obj.BBox[2] * float64(width))
		h := int(obj.BBox[3] * float64(height))

		croppedImg := image.NewRGBA(image.Rect(0, 0, w, h))

		draw.Draw(croppedImg, croppedImg.Bounds(), img, image.Point{X: x, Y: y}, draw.Src)

		centroids, err := prominentcolor.Kmeans(croppedImg)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(centroids) == 0 {
			continue
		}

		newObject := obj.Copy()
		newObject.Color = centroids[0].AsString()
		newObjects = append(newObjects, newObject)
	}

	return Resp{col.objects.Name: newObjects}, nil
}

func (col *Color) GetPorts() (in, out []Port) {
	return []Port{col.frame, col.objects}, []Port{col.objects}
}

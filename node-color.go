package flownodes

import (
	"errors"
	"image"
	"image/draw"

	"github.com/osmoview/flowdata"

	"github.com/EdlinOrg/prominentcolor"
)

func init() {
	RegisterNode("color", &Color{})
}

type Color struct {
	// frame image
	frame Port
	// objects Bbox in relative coordinates format [x,y,w,h]
	objects Port
}

func (col *Color) Init(db DB) (err error) {
	col.initPorts()
	return
}

func (col *Color) initPorts() {
	col.frame = Port{Name: PortFrame, Type: PortFrame}
	col.objects = Port{Name: PortObjects, Type: PortObjects}
}

// getCtxData get frame, objects from Ctx
func (col *Color) getCtxData(c Ctx) (*flowdata.Frame, []*flowdata.Object, error) {
	var frame *flowdata.Frame

	if data, ok := c.GetData(PortFrame); ok {
		frame, ok = data.(*flowdata.Frame)
		if !ok {
			return nil, nil, ErrUnexpectedDataType
		}
	}
	if frame == nil {
		return nil, nil, errors.New("frame not found")
	}

	objects, ok := c.MustGetData(PortObjects).([]*flowdata.Object)
	if !ok {
		return nil, nil, ErrUnexpectedDataType
	}

	return frame, objects, nil
}

func (col *Color) Execute(c Ctx) (out Resp, err error) {
	frame, objects, err := col.getCtxData(c)
	if err != nil {
		return nil, err
	}

	newObjects := make([]*flowdata.Object, 0, len(objects))

	img := frame.GetImage()
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for _, obj := range objects {
		if obj.IsZeroBBox() {
			continue
		}

		x := int(obj.BBox[0] * float64(width))
		y := int(obj.BBox[1] * float64(height))
		w := int(obj.BBox[2] * float64(width))
		h := int(obj.BBox[3] * float64(height))

		cropped := image.NewRGBA(image.Rect(0, 0, w, h))

		draw.Draw(cropped, cropped.Bounds(), img, image.Point{X: x, Y: y}, draw.Src)

		centroids, err := prominentcolor.Kmeans(cropped)
		if err != nil {
			return nil, err
		}

		if len(centroids) == 0 {
			return nil, errors.New("color cropped image not found")
		}

		newObject := obj.Copy()
		newObject.Color = centroids[0].AsString()
		newObjects = append(newObjects, newObject)
	}

	return Resp{col.objects.Name: newObjects}, nil
}

func (col *Color) GetPorts() (in, out []Port) {
	col.initPorts()
	return []Port{col.frame, col.objects}, []Port{col.objects}
}

func (col *Color) GetForm() flowdata.Form {
	return flowdata.Form{}
}

func (col *Color) GetGroup() string {
	return GroupObjects
}

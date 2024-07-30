package flownodes

import (
	"errors"
	"fmt"
	"image"
	"image/draw"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("crop", &Crop{})
}

type Crop struct {
	frame  Port
	region Port
}

func (cr *Crop) Init(db DB) (err error) {
	cr.initPorts()
	return
}

func (cr *Crop) initPorts() {
	cr.frame = Port{Name: PortFrame, Type: PortFrame, Required: true}
	cr.region = Port{Name: PortRegion, Type: PortRegion, Required: true}
}

// getCtxData get frame, region from Ctx
func (cr *Crop) getCtxData(c Ctx) (*flowdata.Frame, *flowdata.Region, error) {
	var frame *flowdata.Frame
	var region *flowdata.Region

	if data, ok := c.GetData(PortFrame); ok {
		frame, ok = data.(*flowdata.Frame)
		if !ok {
			return nil, nil, ErrUnexpectedDataType
		}
	}
	if frame == nil {
		return nil, nil, errors.New("frame not found")
	}

	if data, ok := c.GetData(PortRegion); ok {
		region, ok = data.(*flowdata.Region)
		if !ok {
			return nil, nil, ErrUnexpectedDataType
		}
	}
	if region == nil {
		return nil, nil, errors.New("region not found")
	}

	return frame, region, nil
}

func (cr *Crop) Execute(c Ctx) (out Resp, err error) {
	frame, region, err := cr.getCtxData(c)
	if err != nil {
		return nil, err
	}

	img := frame.GetImage()
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	bbox, ok := region.BBox()
	if !ok {
		return nil, fmt.Errorf("%s is invalid region", region.Name)
	}

	x := int(bbox[0] * float64(width))
	y := int(bbox[1] * float64(height))
	w := int(bbox[2] * float64(width))
	h := int(bbox[3] * float64(height))

	cropped := image.NewRGBA(image.Rect(0, 0, w, h))

	draw.Draw(cropped, cropped.Bounds(), img, image.Point{X: x, Y: y}, draw.Src)

	newFrame, err := flowdata.NewFrameImage(cropped)
	if err != nil {
		return nil, err
	}

	return Resp{cr.frame.Name: newFrame}, nil
}

func (cr *Crop) GetPorts() (in, out []Port) {
	cr.initPorts()
	return []Port{cr.frame, cr.region}, []Port{cr.frame}
}

func (cr *Crop) GetForm() flowdata.Form {
	return flowdata.Form{}
}

func (cr *Crop) GetGroup() string {
	return GroupFrames
}

package flownodes

import (
	"errors"
	"image"
	"image/draw"

	"github.com/osmoview/flowdata"
)

type Crop struct {
	frame  Port
	region Port
}

func (cr *Crop) Init(db DB) (err error) {
	cr.frame = Port{Name: PortFrame, Type: PortFrame}
	cr.region = Port{Name: PortRegion, Type: PortRegion}

	return
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

	bbox := region.BBox()
	if region.IsZeroBBox(bbox) {
		return nil, errors.New("checked bbox is zero")
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
	return []Port{cr.frame, cr.region}, []Port{cr.frame}
}

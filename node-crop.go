package flownodes

import (
	"errors"
	"github.com/osmoview/flowdata"
	"image"
	"image/draw"

	// blank import to support some image formats
	_ "image/jpeg"
	_ "image/png"
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

func (cr *Crop) Execute(c Ctx) (out Resp, err error) {
	var frame *flowdata.Frame
	var region *flowdata.Region

	if data, ok := c.GetData(PortFrame); ok {
		frame, ok = data.(*flowdata.Frame)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
	}
	if frame == nil {
		return nil, errors.New("frame not found")
	}

	if data, ok := c.GetData(PortRegion); ok {
		region, ok = data.(*flowdata.Region)
		if !ok {
			return nil, ErrUnexpectedDataType
		}
	}
	if region == nil {
		return nil, errors.New("region not found")
	}

	img := frame.GetImage()
	minX, minY, maxX, maxY := region.BBox()

	rect := image.Rect(int(minX), int(minY), int(maxX), int(maxY))

	croppedImg := image.NewRGBA(rect.Bounds())

	draw.Draw(croppedImg, rect.Bounds(), img, rect.Min, draw.Src)

	newFrame, err := flowdata.NewFrameImage(croppedImg)
	if err != nil {
		return nil, err
	}

	return Resp{cr.frame.Name: newFrame}, nil
}

func (cr *Crop) GetPorts() (in, out []Port) {
	return []Port{cr.frame, cr.region}, []Port{cr.frame}
}

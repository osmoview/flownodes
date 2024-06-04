package flownodes

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	// blank import to support some image formats
	_ "image/jpeg"
	_ "image/png"

	"github.com/osmoview/flowdata"
)

const (
	formatJPEG = "jpeg"
	formatPNG  = "png"
)

var (
	ErrInvalidImage      = errors.New("invalid image")
	ErrUnsupportedFormat = errors.New("unsupported image format")

	jpegEndSignature = []byte{0xff, 0xd9}
)

type Crop struct {
	// frame image
	frame flowdata.Frame
	// region
	region [][2]float64

	port Port
}

func (cr *Crop) Init(db DB) error {
	cr.port = Port{Name: PortFrame, Type: PortFrame}

	return nil
}

func (cr *Crop) Execute(c Ctx) (out Resp, err error) {
	img := cr.frame.GetImage()

	format, err := cr.getImageFormat(cr.frame.GetBytes())
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	minX, minY := width, height
	maxX, maxY := 0, 0

	for _, point := range cr.region {
		x := int(point[0] * float64(width))
		y := int(point[1] * float64(height))

		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	if minX < 0 || minY < 0 || maxX > width || maxY > height {
		return nil, ErrCoordinatesOutsideImage
	}

	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(minX, minY, maxX, maxY))

	cr.frame.Image = croppedImg
	cr.frame.Info = image.Config{
		ColorModel: croppedImg.ColorModel(),
		Width:      croppedImg.Bounds().Dx(),
		Height:     croppedImg.Bounds().Dy(),
	}

	cr.frame.Bytes, err = cr.imageToBytes(croppedImg, format)
	if err != nil {
		return nil, err
	}

	return Resp{cr.port.Name: cr.frame}, nil
}

func (cr *Crop) GetPorts() (in, out []Port) {
	return in, []Port{cr.port}
}

// getImageFormat get image format from the image bytes
func (cr *Crop) getImageFormat(imgBytes []byte) (format string, err error) {
	_, format, err = image.DecodeConfig(bytes.NewReader(imgBytes))
	if err != nil {
		return
	}

	if format == formatJPEG && !bytes.HasSuffix(imgBytes, jpegEndSignature) {
		err = ErrInvalidImage
	}

	return
}

// imageToBytes convert image.Image to bytes
func (cr *Crop) imageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch format {
	case formatJPEG:
		err = jpeg.Encode(&buf, img, nil)
	case formatPNG:
		err = png.Encode(&buf, img)
	default:
		return nil, ErrUnsupportedFormat
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

package flownodes

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	"github.com/osmoview/flowdata"
)

func createTestImageCrop() ([]byte, image.Image) {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 0, 255})
		}
	}

	jpeg.Encode(&buf, img, nil)

	return buf.Bytes(), img
}

func TestCrop_Execute(t *testing.T) {
	imgBytes, img := createTestImageCrop()
	frame := flowdata.Frame{
		Image: img,
		Bytes: imgBytes,
		Info: image.Config{
			ColorModel: img.ColorModel(),
			Width:      img.Bounds().Dx(),
			Height:     img.Bounds().Dy(),
		},
	}

	tests := []struct {
		name    string
		region  [][2]float64
		wantErr bool
	}{
		{
			name: "valid region",
			region: [][2]float64{
				{0.1, 0.1},
				{0.1, 0.5},
				{0.5, 0.1},
				{0.5, 0.5},
			},
			wantErr: false,
		},
		{
			name: "region outside image",
			region: [][2]float64{
				{1.1, 1.1},
				{1.1, 1.5},
				{1.5, 1.1},
				{1.5, 1.5},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cr := &Crop{
				frame:  frame,
				region: test.region,
			}
			cr.Init(nil)

			_, err := cr.Execute(nil)
			if (err != nil) != test.wantErr {
				t.Errorf("Crop.Execute() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

package frames

import (
	"errors"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

// Crop frame by regions and objects
type Crop struct {
	Regions []string `yaml:"regions"`
	Expand  float32  `yaml:"expand"`

	storage flownodes.Storage
}

type CropReq struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
}

type CropResp struct {
	Frames []*flowdata.Frame
}

func (node *Crop) Init(storage flownodes.Storage) {
	node.storage = storage

	if node.Expand <= 0 {
		node.Expand = 1
	}
}

func (node *Crop) Run(req CropReq, c flownodes.Ctx) (resp CropResp, err error) {
	if req.Frame == nil {
		return resp, errors.New("frame port is required")
	}

	// crop frame by objects bboxes
	for _, o := range req.Objects {
		if !o.BBox.IsValid() {
			continue
		}

		bbox := o.BBox.Scale(node.Expand)

		frame, err := req.Frame.Crop(bbox)
		if err != nil {
			return resp, err
		}

		resp.Frames = append(resp.Frames, frame)
	}

	// crop frame by regions
	// load camera from the storage
	ds, err := node.storage.GetDataSource(c.GetSourceInfo().ID)
	if err != nil {
		return resp, err
	}

	for _, regionName := range node.Regions {
		// get required regions from the camera, and crop frame by it
		for _, r := range ds.GetRegions(regionName) {
			if bbox, ok := r.BBox(); ok && bbox.IsValid() {
				bbox = bbox.Scale(node.Expand)

				frame, err := req.Frame.Crop(bbox)
				if err != nil {
					return resp, err
				}

				resp.Frames = append(resp.Frames, frame)
			}
		}
	}

	return
}

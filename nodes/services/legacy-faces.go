package services

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type LegacyFaces struct {
	ID string

	service flownodes.RemoteService
}

func (node *LegacyFaces) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	return
}

type LegacyFacesReq struct {
	Frame *flowdata.Frame
}

type LegacyFacesResp struct {
	Objects []*flowdata.Object
}

func (node *LegacyFaces) Run(req LegacyFacesReq) (resp LegacyFacesResp, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return resp, nil
	}

	reqData := map[string]any{"frame": req.Frame}
	var faces []LegacyFace
	if err := node.service.Upload("/cvpredict", reqData, &faces); err != nil {
		return resp, err
	}

	for _, face := range faces {
		resp.Objects = append(resp.Objects, face.toObject(req.Frame.Info.Width, req.Frame.Info.Height))
	}

	return
}

type LegacyFace struct {
	Prob   float64
	BBox   [4]int
	Vector flowdata.Vector
}

func (face LegacyFace) toObject(width, height int) *flowdata.Object {
	x := float32(face.BBox[0]) / float32(width)
	y := float32(face.BBox[1]) / float32(height)
	w := float32(face.BBox[2]) / float32(width)
	h := float32(face.BBox[3]) / float32(height)

	return &flowdata.Object{
		Type:   "face",
		Prob:   face.Prob,
		BBox:   [4]float32{x, y, w, h},
		Vector: face.Vector,
	}
}

package nodes

import (
	"fmt"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type Detector struct {
	ID   string
	Path string

	service flownodes.RemoteService
}

func (node *Detector) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	return
}

type DetectorReq struct {
	Frame  *flowdata.Frame
	Frames []*flowdata.Frame
}

func (req *DetectorReq) DataMap() (data map[string]any) {
	// initialize a frames slice to avoid null in the request
	frames := []*flowdata.Frame{}

	// if one frame exist, append it
	if req.Frame != nil {
		frames = append(frames, req.Frame)
	}

	// preapre data request with all frames
	data = map[string]any{
		"frame": append(frames, req.Frames...),
	}

	return
}

func (req *DetectorReq) GetFrame(id string) (*flowdata.Frame, error) {
	if req.Frame != nil && req.Frame.ID == id {
		return req.Frame, nil
	}

	for _, f := range req.Frames {
		if f.ID == id {
			return f, nil
		}
	}

	return nil, fmt.Errorf("frame by id '%s' not exist", id)
}

type responseWithBatch struct {
	Predicts []*flowdata.Object
	Batch    map[string][]*flowdata.Object
}

type DetectorResp struct {
	Objects []*flowdata.Object
}

func (node *Detector) Run(req DetectorReq) (resp DetectorResp, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return resp, nil
	}

	var r responseWithBatch
	if err = node.service.Post(node.Path, req.DataMap(), &r); err != nil {
		return resp, err
	}

	resp.Objects = append(resp.Objects, r.Predicts...)

	for frameid, objects := range r.Batch {
		f, err := req.GetFrame(frameid)
		if err != nil {
			return resp, err
		}

		// frame.BBox valid for sub images
		if f.BBox.IsValid() {
			for _, o := range objects {
				o.BBox = o.BBox.RestoreSize(f.BBox)
			}
		}

		resp.Objects = append(resp.Objects, objects...)
	}

	return
}

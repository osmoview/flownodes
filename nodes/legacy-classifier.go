package nodes

import (
	"encoding/json"
	"errors"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type LegacyClassifier struct {
	ID   string
	Path string

	service flownodes.RemoteService
}

func (node *LegacyClassifier) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	return
}

type LegacyClassifierReq struct {
	Frame   *flowdata.Frame
	Objects []*flowdata.Object
}

type LegacyClassifierResp struct {
	Objects []*flowdata.Object
}

func (node *LegacyClassifier) Run(req LegacyClassifierReq) (resp LegacyClassifierResp, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return resp, nil
	}

	if req.Frame == nil {
		return resp, errors.New("frame is required")
	}

	for _, o := range req.Objects {
		if !o.BBox.IsValid() {
			continue
		}

		subFrame, err := req.Frame.Crop(o.BBox)

		r, err := node.sendFrame(subFrame)
		if err != nil {
			return resp, err
		}

		node.fixBBoxes(subFrame, &r)
		o.Merge(&r)

		resp.Objects = append(resp.Objects, o)
	}

	return
}

type responseWithPredicts struct {
	Predicts flowdata.Object
}

func (node *LegacyClassifier) sendFrame(frame *flowdata.Frame) (obj flowdata.Object, err error) {
	var rawResp json.RawMessage
	if err := node.service.Upload(node.Path, map[string]any{"frame": frame}, &rawResp); err != nil {
		return obj, err
	}

	// try to parse response with predicts and batch fields
	{
		var r responseWithPredicts
		if err = json.Unmarshal(rawResp, &r); err == nil {
			return r.Predicts, nil
		}
	}

	// try to parse response to the array of objects
	{
		if err = json.Unmarshal(rawResp, &obj); err == nil {
			return
		}
	}

	return obj, errors.New("unexpected data layout")
}

func (node *LegacyClassifier) fixBBoxes(frame *flowdata.Frame, objects ...*flowdata.Object) {
	for _, o := range objects {
		if !o.BBox.IsValid() {
			continue
		}

		o.BBox = o.BBox.RestoreSize(frame.BBox)

		if len(o.Objects) > 0 {
			node.fixBBoxes(frame, o.Objects...)
		}
	}
}

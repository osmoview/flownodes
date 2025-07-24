package services

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type LegacyClassifier struct {
	ID   string
	Path string

	service flownodes.RemoteService
	data    flownodes.Data
}

func (node *LegacyClassifier) Init(storage flownodes.Storage, data flownodes.Data) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	node.data = data

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

func (node *LegacyClassifier) sendFrame(frame *flowdata.Frame) (obj flowdata.Object, err error) {
	path := node.Path + "?cam_id=" + frame.DataSource.ID
	data := map[string]any{
		"frame":      frame,
		"datasource": frame.DataSource,
		"node_data":  node.data,
	}

	var rawResp json.RawMessage
	if err := node.service.Upload(path, data, &rawResp); err != nil {
		return obj, err
	}

	// check if response has predicts, then try to parse with it
	if bytes.HasPrefix([]byte(rawResp), []byte(`{"predicts":`)) {
		var r struct{ Predicts flowdata.Object }
		if err = json.Unmarshal(rawResp, &r); err == nil {
			return r.Predicts, nil
		}
	}

	err = json.Unmarshal(rawResp, &obj)
	return
}

func (node *LegacyClassifier) fixBBoxes(frame *flowdata.Frame, objects ...*flowdata.Object) {
	var attributeIndex float32
	for _, o := range objects {
		if !o.BBox.IsValid() {
			o.BBox[1] += attributeIndex * 0.1
			attributeIndex++
		}

		o.BBox = o.BBox.RestoreSize(frame.BBox)

		if len(o.Objects) > 0 {
			node.fixBBoxes(frame, o.Objects...)
		}
		if len(o.Items) > 0 {
			node.fixBBoxes(frame, o.Items...)
		}
	}
}

package services

import (
	"encoding/json"
	"errors"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type LegacyDetector struct {
	ID   string
	Path string

	service flownodes.RemoteService
	form    flownodes.Data
}

func (node *LegacyDetector) Init(storage flownodes.Storage, form flownodes.Data) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	node.form = form

	return
}

type LegacyDetectorReq struct {
	Frame  *flowdata.Frame
	Frames []*flowdata.Frame
}

type LegacyDetectorResp struct {
	Objects []*flowdata.Object
}

func (node *LegacyDetector) Run(req LegacyDetectorReq, input flownodes.Data) (resp LegacyDetectorResp, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return resp, nil
	}

	if req.Frame != nil {
		resp.Objects, err = node.sendFrame(req.Frame, input)
		if err != nil {
			return
		}

		node.fixBBoxes(req.Frame, resp.Objects)
	}

	for _, frame := range req.Frames {
		objects, err := node.sendFrame(frame, input)
		if err != nil {
			return resp, err
		}

		node.fixBBoxes(frame, objects)

		resp.Objects = append(resp.Objects, objects...)
	}

	return
}

func (node *LegacyDetector) sendFrame(frame *flowdata.Frame, input flownodes.Data) (objects []*flowdata.Object, err error) {
	path := node.Path + "?cam_id=" + frame.DataSource.ID
	data := map[string]any{
		"frame":      frame,
		"datasource": frame.DataSource,
		"node_data":  node.mergeData(input),
	}

	var rawResp json.RawMessage
	if err := node.service.Upload(path, data, &rawResp); err != nil {
		return nil, err
	}

	// try to parse response with predicts and batch fields
	{
		var r struct{ Predicts []*flowdata.Object }
		if err = json.Unmarshal(rawResp, &r); err == nil {
			return r.Predicts, nil
		}
	}

	// try to parse response to the array of objects
	{
		if err = json.Unmarshal(rawResp, &objects); err == nil {
			return
		}
	}

	return nil, errors.New("unexpected data layout")
}

func (node *LegacyDetector) fixBBoxes(frame *flowdata.Frame, objects []*flowdata.Object) {
	if !frame.BBox.IsValid() {
		return
	}

	for _, o := range objects {
		o.BBox = o.BBox.RestoreSize(frame.BBox)

		if len(o.Objects) > 0 {
			node.fixBBoxes(frame, o.Objects)
		}
	}
}

func (node *LegacyDetector) mergeData(input flownodes.Data) flownodes.Data {
	data := flownodes.Data{}
	for k, v := range node.form {
		data[k] = v
		if v, ok := input[k]; ok {
			data[k] = v
		}
	}

	return data
}

package nodes

import (
	"encoding/json"
	"errors"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type LegacyDetector struct {
	ID string

	service flownodes.RemoteService
}

func (node *LegacyDetector) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	return
}

type responseWithPredicts struct {
	Predicts []*flowdata.Object
}

type LegacyDetectorReq struct {
	Frame  *flowdata.Frame
	Frames []*flowdata.Frame
}

type LegacyDetectorResp struct {
	Objects []*flowdata.Object
}

func (node *LegacyDetector) Run(req LegacyDetectorReq) (resp LegacyDetectorResp, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return resp, nil
	}

	if req.Frame != nil {
		resp.Objects, err = node.sendFrame(req.Frame)
		if err != nil {
			return
		}
	}

	for _, frame := range req.Frames {
		objects, err := node.sendFrame(frame)
		if err != nil {
			return resp, err
		}
		if frame.BBox.IsValid() {
			for _, o := range objects {
				o.BBox = o.BBox.RestoreSize(frame.BBox)
			}
		}
		resp.Objects = append(resp.Objects, objects...)
	}

	return
}

func (node *LegacyDetector) sendFrame(frame *flowdata.Frame) (objects []*flowdata.Object, err error) {
	reqData := map[string]any{"frame": frame}

	// send request
	var rawResp json.RawMessage
	if err = node.service.Send(reqData, &rawResp); err != nil {
		return nil, err
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
		if err = json.Unmarshal(rawResp, &objects); err == nil {
			return
		}
	}

	return nil, errors.New("unexpected data layout")
}

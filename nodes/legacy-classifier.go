package nodes

import (
	"errors"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type LegacyClassifier struct {
	ID string

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

		if err := node.sendFrame(req.Frame, o); err != nil {
			return resp, err
		}
	}

	return
}

func (node *LegacyClassifier) sendFrame(frame *flowdata.Frame, o *flowdata.Object) (err error) {
	frame, err = frame.Crop(o.BBox)
	if err != nil {
		return err
	}

	reqData := map[string]any{"frame": frame}

	// send request
	var resp flowdata.Object
	if err := node.service.Send(reqData, &resp); err != nil {
		return err
	}

	node.fixBBoxes(frame, &resp)

	// if classifier change object type, change it manually for the object as o.Merge do not change type
	if resp.Type == "" {
		o.Type = resp.Type
	}

	o.Merge(&resp)

	return nil
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

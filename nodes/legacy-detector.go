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

type respWithPredicts struct {
	Predicts []*flowdata.Object
}

func (node *LegacyDetector) Run(c flownodes.Ctx) (resp flownodes.Data, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return nil, nil
	}

	var rawResp json.RawMessage
	if err = node.service.Send(c.DataMap(), &rawResp); err != nil {
		return resp, err
	}
	{
		var resp respWithPredicts
		if err = json.Unmarshal(rawResp, &resp); err == nil {
			return flownodes.Data{flownodes.PortObjects: resp.Predicts}, nil
		}
	}

	{
		var list []*flowdata.Object
		if err = json.Unmarshal(rawResp, &list); err == nil {
			return flownodes.Data{flownodes.PortObjects: list}, nil
		}
	}

	return resp, errors.New("unexpected data layout")
}

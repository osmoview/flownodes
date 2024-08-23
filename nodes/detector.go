package nodes

import (
	"github.com/osmoview/flownodes"
)

type Detector struct {
	ID string

	service flownodes.RemoteService
}

func (node *Detector) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	return
}

func (node *Detector) Run(c flownodes.Ctx) (resp flownodes.Data, err error) {
	// WARN: i do not know what to do if there is no service
	if node.service == nil {
		return nil, nil
	}

	err = node.service.Send(c.DataMap(), &resp)
	return
}

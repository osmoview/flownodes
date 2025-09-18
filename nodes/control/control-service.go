package control

import (
	"errors"

	"github.com/osmoview/flownodes"
)

type ControlService struct {
	RemoteID string
	service  flownodes.RemoteService
}

func (node *ControlService) Init(storage flownodes.Storage) {
	node.service, _ = storage.GetService(node.RemoteID)
}

func (node *ControlService) Run(c flownodes.Ctx) (err error) {
	if node.service == nil {
		return errors.New("service not exist")
	}

	return node.service.Update(c.DataMap())
}

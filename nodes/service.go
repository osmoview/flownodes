package nodes

import (
	"github.com/osmoview/flownodes"
)

type Service struct {
	ID string

	service flownodes.RemoteService
}

// Init should load client to the service by id
func (node *Service) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}

	return
}

// Run send incoming data(map[string]any) to the remote service
// and should receive response as a compatibe with map[string]any format
// WARN: string, arrays and other responses not supported
func (node *Service) Run(c flownodes.Ctx) (resp flownodes.Data, err error) {
	if node.service == nil {
		return nil, nil
	}

	err = node.service.Send(c.DataMap(), &resp)
	return
}

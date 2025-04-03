package nodes

import (
	"errors"

	"github.com/osmoview/flownodes"
)

type RPC struct {
	ID string

	service flownodes.RemoteService
}

func (node *RPC) Init(data flownodes.Data, storage flownodes.Storage) (err error) {
	if node.ID == "" {
		return errors.New("no rpc service ID")
	}

	node.service, err = storage.GetService(node.ID)
	if err != nil {
		return err
	}

	// if len(data) > 0 {
	// 	return node.service.Send(map[string]any{"init": data}, nil, nil)
	// }

	return nil
}

func (node *RPC) Run(data flownodes.Data) (resp flownodes.Data, err error) {
	if node.service == nil {
		return nil, errors.New("no rpc service")
	}

	err = node.service.Post("", map[string]any{"run": data}, &resp)
	return
}

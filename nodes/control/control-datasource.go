package control

import (
	"errors"

	"github.com/osmoview/flownodes"
)

type ControlDataSource struct {
	RemoteID   string
	datasource flownodes.DataSource
}

func (node *ControlDataSource) Init(storage flownodes.Storage) (err error) {
	node.datasource, err = storage.GetDataSource(node.RemoteID)
	return
}

func (node *ControlDataSource) Run(c flownodes.Ctx) (err error) {
	if node.datasource == nil {
		return errors.New("datasource not exist")
	}

	return node.datasource.Update(c.DataMap())
}

package nodes

import (
	"errors"

	"github.com/osmoview/flownodes"
)

type DBSelect struct {
	Table string
	db    flownodes.Database
}

type DBSelectData struct {
	Record flownodes.Data
}

func (node *DBSelect) Init(db flownodes.Database) error {
	if node.Table == "" {
		return errors.New("table is required")
	}

	node.db = db

	return nil
}

func (node *DBSelect) Run(data flownodes.Data) (resp DBSelectData, err error) {
	table, err := node.db.Table(flownodes.TableScheme{Table: node.Table})
	if err != nil {
		return resp, err
	}

	record, err := table.SelectOne(data)
	if err != nil {
		return resp, err
	}

	return DBSelectData{Record: record}, nil
}

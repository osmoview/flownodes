package nodes

import (
	"errors"

	"github.com/osmoview/flownodes"
)

type DBInsert struct {
	Table string
	db    flownodes.Database
}

type DBInsertData struct {
	Record flownodes.Data
}

func (node *DBInsert) Init(db flownodes.Database) error {
	if node.Table == "" {
		return errors.New("table is required")
	}

	return nil
}

func (node *DBInsert) Run(data flownodes.Data) (DBInsertData, error) {
	table, err := node.db.Table(flownodes.TableScheme{Table: node.Table})
	if err != nil {
		return DBInsertData{}, err
	}

	if err := table.Insert(data); err != nil {
		return DBInsertData{}, err
	}

	return DBInsertData{Record: data}, nil
}

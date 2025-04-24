package nodes

import (
	"log"

	"github.com/osmoview/flownodes"
)

type DBRecordExample struct {
	Table string
	Value int

	db flownodes.Database
}

func (node *DBRecordExample) Init(db flownodes.Database) {
	node.db = db

	if node.Table == "" {
		node.Table = "test-table"
	}
}

func (node *DBRecordExample) Run(data flownodes.Data) (resp flownodes.Data, err error) {
	tm := flownodes.TableScheme{
		Table: node.Table,
		Columns: []flownodes.TableColumn{
			{Name: "id", Type: "bigint NOT NULL AUTO_INCREMENT"},
			{Name: "name", Type: "text"},
			{Name: "value", Type: "bigint"},
		},
	}

	table, err := node.db.Table(tm)
	if err != nil {
		return resp, err
	}

	if err := table.Insert(data); err != nil {
		return data, err
	}

	log.Printf("%+v", data)

	rec, err := table.SelectOne(map[string]any{"name": "World"})
	if err != nil {
		return resp, err
	}

	log.Printf("%+v", rec)

	return data, nil
}

type RecordExample struct {
	Table string `gorm:"-"`

	Name  string `gorm:"index"`
	Value int
}

func (r RecordExample) TableName() string {
	return r.Table
}

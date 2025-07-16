package simple

import (
	"log"

	"github.com/osmoview/flownodes"
)

type TestNode struct {
}

func (node *TestNode) Run(data flownodes.Data) flownodes.Data {
	log.Printf("%+v", data)
	return data
}

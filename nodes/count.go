package nodes

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type Count struct {
	Count int `yaml:"count"`
}

type CountReq struct {
	Objects []*flowdata.Object
}

func (node *Count) Run(req CountReq) flownodes.Data {
	resp := flownodes.Data{}
	if len(req.Objects) > node.Count {
		resp["more"] = req.Objects
	} else if len(req.Objects) == node.Count {
		resp["equal"] = req.Objects
	} else {
		resp["less"] = req.Objects
	}

	return resp
}

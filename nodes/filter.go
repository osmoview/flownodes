package nodes

import (
	"github.com/osmoview/flowdata"
)

const FilterAny = "*"

type Filter struct {
	Probs map[string]float64 `yaml:"probs"`
}

type FilterData struct {
	Objects []*flowdata.Object
}

func (f *Filter) Run(req FilterData) (resp FilterData) {
	resp.Objects = make([]*flowdata.Object, 0, len(req.Objects))

	for _, o := range req.Objects {
		prob, ok := f.Probs[o.Type]
		if !ok {
			prob, ok = f.Probs[FilterAny]
		}

		if !ok || prob > o.Prob {
			continue
		}

		resp.Objects = append(resp.Objects, o)
	}

	return
}

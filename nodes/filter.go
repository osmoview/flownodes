package nodes

import (
	"github.com/osmoview/flowdata"
)

const FilterAny = "*"

type Filter struct {
	Prob    float64  `yaml:"prob" json:"prob"`
	Classes []string `yaml:"classes" json:"classes"`
}

type FilterData struct {
	Objects []*flowdata.Object `json:"objects"`
}

func (f *Filter) Run(req FilterData) (resp FilterData) {
	resp.Objects = make([]*flowdata.Object, 0, len(req.Objects))

	for _, o := range req.Objects {
		if o.Prob < f.Prob || !f.classContains(o.Type) {
			continue
		}

		resp.Objects = append(resp.Objects, o)
	}

	return
}

func (f *Filter) classContains(s string) bool {
	if len(f.Classes) == 0 {
		return true
	}

	for _, c := range f.Classes {
		if c == s {
			return true
		}
	}

	return false
}

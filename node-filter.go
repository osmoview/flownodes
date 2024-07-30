package flownodes

import (
	"fmt"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("filter", &Filter{})
}

const FilterAny = "*"

type Filter struct {
	Probs map[string]float64 `yaml:"probs"`

	In  Port `yaml:"-"`
	Out Port `yaml:"-"`
}

func (f *Filter) Init(db DB) error {
	f.Probs = make(map[string]float64)
	f.initPorts()

	return nil
}

func (f *Filter) initPorts() {
	f.In = Port{Name: "objects", Type: PortObjects}
	f.Out = Port{Name: "objects", Type: PortObjects}
}

func (f *Filter) Execute(c Ctx) (out Resp, err error) {
	var probs = f.Probs

	list, ok := c.MustGetData("objects").([]*flowdata.Object)
	if !ok {
		err = fmt.Errorf("input data not an objects type")
		return
	}

	resp := make([]*flowdata.Object, 0, 0)
	for _, o := range list {
		prob, ok := probs[o.Type]
		if !ok {
			prob, ok = probs[FilterAny]
		}
		if !ok || prob > o.Prob {
			continue
		}

		resp = append(resp, o)

	}

	return Resp{"objects": resp}, nil
}

func (f *Filter) GetPorts() (in []Port, out []Port) {
	f.initPorts()
	return []Port{f.In}, []Port{f.Out}
}

func (f *Filter) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:    "probs",
			Label:   flowdata.L10n{"en": "Threshold", "ru": "Порог"},
			Type:    flowdata.FieldTypeMapNumbers,
			Default: map[string]float64{"*": 0.35},
			Help:    flowdata.L10n{"en": "Filter threshold"},
		},
	}
}

func (f *Filter) GetGroup() string {
	return GroupObjects
}

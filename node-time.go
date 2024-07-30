package flownodes

import (
	"time"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("time", &Time{})
	SourceNodes = append(SourceNodes, "time")
}

type Time struct {
	Layout string `yaml:"layout"`

	output Port
}

func (r *Time) Init(db DB) error {
	if r.Layout == "" {
		r.Layout = "2006-01-02_15:04:05"
	}

	r.initPorts()
	return nil
}

func (r *Time) initPorts() {
	r.output = Port{Name: "time", Type: PortString}
}

func (r *Time) Execute(c Ctx) (out Resp, err error) {
	s := time.Now().Format(r.Layout)

	return Resp{r.output.Name: s}, nil
}

func (r *Time) GetPorts() (in, out []Port) {
	r.initPorts()
	return []Port{}, []Port{r.output}
}

func (r *Time) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:  "layout",
			Label: flowdata.L10n{"en": "Layout", "ru": "Формат"},
			Help:  flowdata.L10n{"en": "Time Format Layout", "ru": "Шаблон форматирования времени"},
			Level: flowdata.FieldLevelMaster,
		},
	}
}

func (t *Time) GetGroup() string {
	return GroupSimple
}

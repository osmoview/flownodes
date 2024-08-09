package flownodes

import (
	"errors"

	"github.com/tidwall/gjson"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("gjson", &GJSON{})
}

type GJSON struct {
	Expr string `yaml:"expr"`

	in  Port
	out Port
}

func (e *GJSON) Init(db DB) error {
	e.initPorts()
	return nil
}

func (e *GJSON) initPorts() {
	e.in = Port{Name: "data", Type: PortAny}
	e.out = Port{Name: "data", Type: PortAny}
}

func (e *GJSON) getInputData(c Ctx) (string, error) {
	m := c.MustGetData("data")

	if b, ok := m.([]byte); ok {
		return string(b), nil
	}

	if s, ok := m.(string); ok {
		return s, nil
	}

	return "", errors.New("unexpected data type")
}

func (e *GJSON) Execute(c Ctx) (Resp, error) {
	s, err := e.getInputData(c)
	if err != nil {
		return nil, err
	}

	log.Debug(s)

	result := gjson.Get(s, e.Expr)

	log.Debug(result)

	return Resp{e.out.Name: result.Value()}, nil
}

func (e *GJSON) GetPorts() (in, out []Port) {
	e.initPorts()
	return []Port{e.in}, []Port{e.out}
}

func (e *GJSON) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:  "expr",
			Help:  flowdata.L10n{"en": "JSON path expresion"},
			Label: flowdata.L10n{"en": "Path Expression", "ru": "Путь выражение"},
			Level: flowdata.FieldLevelMaster,
		},
		{
			Name:    "url",
			Type:    flowdata.FieldTypeURL,
			Default: "https://github.com/tidwall/gjson",
			Level:   flowdata.FieldLevelMaster,
		},
	}
}

func (e *GJSON) GetGroup() string {
	return GroupSimple
}

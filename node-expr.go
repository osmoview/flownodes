package flownodes

import (
	"github.com/expr-lang/expr"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("expr", &Expression{})
}

// Expression node use package 'github.com/expr-lang/expr' to implement simple custom logic.
// At the input, in the "data" port, it takes a value of any type,
// then presents the data in the form of a 'Map' and transmits it to the "github.com/expr-lang/expr" packet
type Expression struct {
	// Expr string, for example:
	// 'data > 50', to check that a some number is greater than the required value
	Expr string `yaml:"expr"`

	in  Port
	out Port
}

func (e *Expression) Init(db DB) error {
	e.initPorts()
	return nil
}

func (e *Expression) initPorts() {
	e.in = Port{Name: "*", Type: PortMap}
	e.out = Port{Name: "data", Type: PortAny}
}

func (e *Expression) Execute(c Ctx) (Resp, error) {
	m := c.DataMap()

	program, err := expr.Compile(e.Expr, expr.Env(m))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	out, err := expr.Run(program, m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return Resp{e.out.Name: out}, nil
}

func (e *Expression) GetPorts() (in, out []Port) {
	e.initPorts()
	return []Port{e.in}, []Port{e.out}
}

func (e *Expression) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:  "expr",
			Help:  flowdata.L10n{"en": "Expression string"},
			Label: flowdata.L10n{"en": "Expression", "ru": "Условие"},
			Level: flowdata.FieldLevelMaster,
		},
		{
			Name:    "url",
			Type:    flowdata.FieldTypeURL,
			Default: "https://github.com/expr-lang/expr",
			Level:   flowdata.FieldLevelMaster,
		},
	}
}

func (e *Expression) GetGroup() string {
	return GroupSimple
}

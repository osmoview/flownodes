package flownodes

import (
	"github.com/expr-lang/expr"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("expr", NewExprNode)
}

// Expr node use package 'github.com/expr-lang/expr' to implement simple custom logic.
// At the input, in the "data" port, it takes a value of any type,
// then presents the data in the form of a 'Map' and transmits it to the "github.com/expr-lang/expr" packet
type Expr struct {
	// Expr string, for example:
	// 'data > 50', to check that a some number is greater than the required value
	Expr string `yaml:"expr"`

	in  Port
	out Port
}

func NewExprNode() Node {
	return &Expr{}
}

func (e *Expr) Init(db DB) error {
	e.initPorts()
	return nil
}

func (e *Expr) initPorts() {
	e.in = Port{Name: "*", Type: PortMap}
	e.out = Port{Name: "data", Type: PortAny}
}

func (e *Expr) Execute(c Ctx) (Resp, error) {
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

func (e *Expr) GetPorts() (in, out []Port) {
	e.initPorts()
	return []Port{e.in}, []Port{e.out}
}

func (e *Expr) GetForm() flowdata.Form {
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

func (e *Expr) GetGroup() string {
	return GroupSimple
}

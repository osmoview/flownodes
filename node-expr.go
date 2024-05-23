package flownodes

import (
	"log"

	"github.com/expr-lang/expr"
)

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
	e.in = Port{Name: "*", Type: PortMap}
	e.out = Port{Name: "data", Type: PortAny}
	return nil
}

func (e *Expression) Execute(c Ctx) (Resp, error) {
	m := c.ToMap()

	program, err := expr.Compile(e.Expr, expr.Env(m))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	out, err := expr.Run(program, m)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return Resp{e.out.Name: out}, nil
}

func (e *Expression) GetPorts() (in, out []Port) {
	return []Port{e.in}, []Port{e.out}
}

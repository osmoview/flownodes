package simple

import "github.com/osmoview/flownodes"

type Number struct {
	Number float64
}

func (node *Number) Run() flownodes.Data {
	return flownodes.Data{"number": node.Number}
}

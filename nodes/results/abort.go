package results

import "github.com/osmoview/flownodes"

type Abort struct{}

func (a *Abort) Run(c flownodes.Ctx) {
	c.Abort()
}

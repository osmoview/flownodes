package nodes

import (
	"time"

	"github.com/osmoview/flownodes"
)

type Timer struct {
	Timeout float32
	timeout time.Duration

	t time.Time
}

func (node *Timer) Init() {
	node.timeout = time.Duration(node.Timeout * float32(time.Second))
}

func (node *Timer) Run(c flownodes.Ctx) flownodes.Data {
	if t := c.GetTime(); t.Before(node.t) {
		node.t = t
		return c.DataMap()
	}

	return map[string]any{}
}

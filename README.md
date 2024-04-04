# FlowNodes

WIP!!!

This repository contains the standard `nodes` used in analytics(`flow`).

Each node must correspond to the interface `Node`:

```go
type Node interface {
	// initialize node, set default values, etc
	Init()

	// this function is called every time incoming data(`Ctx``) is received,
	// and should returns response(Resp:map[string]any),
	// where the key is a port name, and the value contains the data
	Execute(c Ctx) (out Resp, err error)

	// get input and output node ports
	GetPorts() (in, out []Port)
}

```

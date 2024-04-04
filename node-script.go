package nodes

type Script struct {
	ID string

	in  []Port
	out []Port
}

func (s *Script) Init() {
	// TODO: init script, resolve ports by public functions
}

func (s *Script) Execute(c Ctx) (out Resp, err error) {
	// TODO: sync / async execution

	return nil, nil
}

func (s *Script) GetPorts() (in, out []Port) {
	// TODO: return ports based on the script functions
	return s.in, s.out
}

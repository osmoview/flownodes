package flownodes

import "os"

type ReadFile struct {
	Path string `yaml:"path"`

	port Port
}

func (r *ReadFile) Init(db DB) error {
	r.port = Port{Name: "data", Type: PortBytes}
	return nil
}

func (r *ReadFile) Execute(c Ctx) (out Resp, err error) {
	b, err := os.ReadFile(r.Path)
	if err != nil {
		return nil, err
	}

	return Resp{r.port.Name: b}, nil
}

func (r *ReadFile) GetPorts() (in, out []Port) {
	return nil, []Port{r.port}
}

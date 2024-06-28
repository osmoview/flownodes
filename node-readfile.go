package flownodes

import (
	"os"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("readfile", &ReadFile{})
	SourceNodes = append(SourceNodes, "readfile")
}

type ReadFile struct {
	Path string `yaml:"path"`

	port Port
}

func (r *ReadFile) Init(db DB) error {
	r.initPorts()
	return nil
}

func (r *ReadFile) initPorts() {
	r.port = Port{Name: "data", Type: PortBytes}
}

func (r *ReadFile) Execute(c Ctx) (out Resp, err error) {
	b, err := os.ReadFile(r.Path)
	if err != nil {
		return nil, err
	}

	return Resp{r.port.Name: b}, nil
}

func (r *ReadFile) GetPorts() (in, out []Port) {
	r.initPorts()
	return nil, []Port{r.port}
}

func (r *ReadFile) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:  "path",
			Label: flowdata.L10n{"en": "Path", "ru": "Путь к файлу"},
			Help:  flowdata.L10n{"en": "Path to the file to read", "ru": "Путь к файлу для чтения"},
			Level: flowdata.FieldLevelMaster,
		},
	}
}

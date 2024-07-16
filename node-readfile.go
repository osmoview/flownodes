package flownodes

import (
	"os"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/settings"
)

func init() {
	RegisterNode("readfile", &ReadFile{})
	SourceNodes = append(SourceNodes, "readfile")
}

type ReadFile struct {
	Path string `yaml:"path"`

	input  Port
	output Port
}

func (r *ReadFile) Init(db DB) error {
	r.initPorts()
	return nil
}

func (r *ReadFile) initPorts() {
	r.input = Port{Name: "*", Type: PortMap}
	r.output = Port{Name: "data", Type: PortBytes}
}

func (r *ReadFile) Execute(c Ctx) (out Resp, err error) {
	fp, err := settings.Map(c.DataMap()).SubstituteTo(r.Path)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	return Resp{r.output.Name: b}, nil
}

func (r *ReadFile) GetPorts() (in, out []Port) {
	r.initPorts()
	return []Port{r.input}, []Port{r.output}
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

package flownodes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/settings"
	"gopkg.in/yaml.v3"
)

func init() {
	RegisterNode("savefile", NewSaveFileNode)
	SourceNodes = append(SourceNodes, "savefile")
}

type SaveFile struct {
	Path string `yaml:"path"`

	data Port
	args Port
}

func NewSaveFileNode() Node {
	return &SaveFile{}
}

func (r *SaveFile) Init(db DB) error {
	r.initPorts()
	return nil
}

func (r *SaveFile) initPorts() {
	r.data = Port{Name: "data", Type: PortAny, Required: true}
	r.args = Port{Name: "*", Type: PortMap}
}

func (r *SaveFile) Execute(c Ctx) (out Resp, err error) {
	fp, err := settings.Map(c.DataMap()).SubstituteTo(r.Path)
	if err != nil {
		return nil, err
	}

	data, ok := c.GetData("data")
	if !ok {
		return nil, fmt.Errorf("port `data` has no data, nothing to save")
	}

	switch d := data.(type) {
	case *flowdata.Frame:
		return nil, os.WriteFile(fp, d.GetBytes(), 0666)
	case string:
		return nil, os.WriteFile(fp, []byte(d), 0666)
	case []byte:
		return nil, os.WriteFile(fp, d, 0666)
	case int, int64, int32, float32, float64, uint, uint16, uint32, uint64:
		return nil, os.WriteFile(fp, []byte(fmt.Sprint(d)), 0666)
	default:
		return nil, r.saveByExt(fp, d)
	}
}

func (r *SaveFile) saveByExt(fp string, data any) error {
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	switch filepath.Ext(fp) {
	case ".yml", ".yaml":
		e := yaml.NewEncoder(f)
		defer e.Close()
		return e.Encode(data)
	// case ...
	// case ...
	default:
		return json.NewEncoder(f).Encode(data)
	}
}

func (r *SaveFile) GetPorts() (in, out []Port) {
	r.initPorts()
	return []Port{r.data, r.args}, []Port{}
}

func (r *SaveFile) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:  "path",
			Label: flowdata.L10n{"en": "Path", "ru": "Путь к файлу"},
			Help:  flowdata.L10n{"en": "Path to the file where to save the data", "ru": "Путь к файлу куда сохранять данные"},
			Level: flowdata.FieldLevelMaster,
		},
	}
}

func (r *SaveFile) GetGroup() string {
	return GroupFiles
}

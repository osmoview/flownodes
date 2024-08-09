package flownodes

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/settings"
)

const TextNode = "text"

func init() {
	RegisterNode(TextNode, &Text{})
	SourceNodes = append(SourceNodes, TextNode)
}

type Text struct {
	Text string `yaml:"text"`

	in  Port
	out Port
}

func (t *Text) Init(db DB) error {
	t.initPorts()
	return nil
}

func (t *Text) initPorts() {
	t.in = Port{Name: "*", Type: PortMap}
	t.out = Port{Name: "text", Type: PortString}
}

func (t *Text) Execute(c Ctx) (out Resp, err error) {
	// substite data to the text line
	s, err := settings.Map(c.DataMap()).SubstituteTo(t.Text)
	if err != nil {
		return nil, err
	}

	return Resp{t.out.Name: s}, nil
}

func (t *Text) GetPorts() (in, out []Port) {
	t.initPorts()
	return []Port{t.in}, []Port{t.out}
}

func (t *Text) GetForm() flowdata.Form {
	return flowdata.Form{
		{Name: "text", Label: flowdata.L10n{"en": "Text", "ru": "Текст"}},
	}
}

func (t *Text) GetGroup() string {
	return GroupSimple
}

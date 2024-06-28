package flownodes

import "github.com/osmoview/flowdata"

const TextNode = "text"

func init() {
	RegisterNode(TextNode, &Text{})
	SourceNodes = append(SourceNodes, TextNode)
}

type Text struct {
	Text string `yaml:"text"`

	port Port
}

func (t *Text) Init(db DB) error {
	t.initPorts()
	return nil
}

func (t *Text) initPorts() {
	t.port = Port{Name: "text", Type: PortString}
}

func (t *Text) Execute(c Ctx) (out Resp, err error) {
	return Resp{t.port.Name: t.Text}, nil
}

func (t *Text) GetPorts() (in, out []Port) {
	t.initPorts()
	return nil, []Port{t.port}
}

func (t *Text) GetForm() flowdata.Form {
	return flowdata.Form{
		{Name: "text", Label: flowdata.L10n{"en": "Text", "ru": "Текст"}},
	}
}

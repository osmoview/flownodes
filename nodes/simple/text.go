package simple

import (
	"bytes"
	"text/template"

	"github.com/osmoview/flownodes"
)

type Text struct {
	Text string
	tmpl *template.Template
}

type TextData struct {
	Text string
}

func (t *Text) Init() (err error) {
	t.tmpl, err = template.New("").Parse(t.Text)
	return
}

func (t *Text) Run(c flownodes.Ctx) (resp TextData, err error) {
	s, err := t.renderTextTemplate(c.DataMap())
	return TextData{Text: s}, err
}

func (t *Text) renderTextTemplate(data map[string]any) (string, error) {
	if t.tmpl != nil {
		var buf bytes.Buffer
		err := t.tmpl.Execute(&buf, data)
		return buf.String(), err
	}

	return t.Text, nil
}

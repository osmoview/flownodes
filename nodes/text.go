package nodes

type Text struct {
	Text string
}

type TextData struct {
	Text string
}

func (t *Text) Run() (resp TextData) {
	return TextData{Text: t.Text}
}

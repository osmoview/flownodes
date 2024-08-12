package flownodes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/osmoview/flowdata"
)

func init() {
	RegisterNode("http", NewHTTPNode)
	SourceNodes = append(SourceNodes, "http")
}

type HTTP struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`

	in       Port
	out      Port
	outError Port
}

func NewHTTPNode() Node {
	return &HTTP{Method: "GET"}
}

func (e *HTTP) Init(db DB) error {
	e.initPorts()
	if e.Method == "" {
		e.Method = "GET"
	}

	return nil
}

func (e *HTTP) initPorts() {
	e.in = Port{Name: "*", Type: PortMap}
	e.out = Port{Name: "data", Type: PortBytes}
	e.outError = Port{Name: "error", Type: PortBytes}
}

func (e *HTTP) Execute(c Ctx) (Resp, error) {
	m := c.MustGetData("data")

	reqdata, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(strings.ToUpper(e.Method), e.URL, bytes.NewReader(reqdata))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respdata, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return Resp{e.outError.Name: respdata}, nil
	}

	return Resp{e.out.Name: respdata}, nil
}

func (e *HTTP) GetPorts() (in, out []Port) {
	e.initPorts()
	return []Port{e.in}, []Port{e.out}
}

func (e *HTTP) GetForm() flowdata.Form {
	return flowdata.Form{
		{
			Name:  "url",
			Label: flowdata.L10n{"en": "URL", "ru": "URL-адрес"},
			Level: flowdata.FieldLevelMaster,
		},
		{
			Name:    "method",
			Label:   flowdata.L10n{"en": "Method", "ru": "Метод"},
			Level:   flowdata.FieldLevelMaster,
			Default: "GET",
		},
	}
}

func (e *HTTP) GetGroup() string {
	return GroupSimple
}

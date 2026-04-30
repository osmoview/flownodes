package services

import (
	"log"
	"net/http"
	"strings"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type Service struct {
	ID          string
	Method      string
	Path        string
	ContentType string `json:"content-type" yaml:"content-type"`

	service flownodes.RemoteService
	form    flownodes.Data
}

// Init should load client to the service by id
func (node *Service) Init(storage flownodes.Storage, form flownodes.Data) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}
	if node.Method == "" {
		node.Method = http.MethodPost
	}

	node.form = form

	return
}

// Run send incoming data(map[string]any) to the remote service
// and should receive response as a compatibe with map[string]any format
// WARN: string, arrays and other responses not supported
func (node *Service) Run(c flownodes.Ctx) (resp flownodes.Data, err error) {
	if node.service == nil {
		log.Println("WARN: no service")
		return nil, nil
	}

	// build URL address
	url, err := node.service.GetAddress(node.Path)
	if err != nil {
		return resp, err
	}

	// create new http request
	req := flowdata.NewRequest()
	req.Method = strings.ToUpper(node.Method)
	req.URL = url
	req.ContentType = node.ContentType

	// set datasource id, flow id, flow run id as http headers
	req.Headers = c.Headers()

	// prepare request data
	data := c.DataMap()
	// set data from the form if it not exists in the input,
	// to prevent overwrite it, as input data more important than form data
	for k, v := range node.form {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}

	if req.Method == http.MethodGet {
		req.Query = data
	} else {
		req.Body = data
	}

	// do request and receive response
	r, err := req.Do()
	if err != nil {
		return resp, err
	}

	if r.StatusCode == 204 {
		return resp, nil
	}

	// parse response as JSON
	return resp, r.JSON(&resp)
}

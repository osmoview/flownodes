package services

import (
	"log"
	"net/http"
	"strings"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type Service struct {
	ID     string
	Method string
	Path   string

	service flownodes.RemoteService
	data    flownodes.Data
}

// Init should load client to the service by id
func (node *Service) Init(storage flownodes.Storage, data flownodes.Data) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}
	if node.Method == "" {
		node.Method = http.MethodPost
	}

	node.data = data

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

	// set datasource id, flow id, flow run id as http headers
	req.Headers = c.Headers()

	// prepare request data
	if req.Method == http.MethodGet {
		req.Query = c.DataMap()
	} else {
		data := c.DataMap()
		data["node_data"] = node.data // substitute node data to the request
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

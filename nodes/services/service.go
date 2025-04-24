package nodes

import (
	"log"
	"net/http"
	"strings"

	"github.com/osmoview/flownodes"
)

type Service struct {
	ID     string
	Method string
	Path   string

	service flownodes.RemoteService
}

// Init should load client to the service by id
func (node *Service) Init(storage flownodes.Storage) (err error) {
	if node.ID != "" {
		node.service, err = storage.GetService(node.ID)
	}
	if node.Method == "" {
		node.Method = http.MethodPost
	}

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

	switch strings.ToUpper(node.Method) {
	case http.MethodPost:
		err = node.service.Post(node.Path, c.DataMap(), &resp)
	case http.MethodGet:
		err = node.service.Get(node.Path, &resp)
	}

	return
}

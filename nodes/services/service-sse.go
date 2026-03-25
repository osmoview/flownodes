package services

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type ServiceSSE struct {
	ID    string
	Flow  string
	Path  string
	Model string

	service flownodes.RemoteService
	flow    flownodes.Flow
}

// Init should load client to the service by id
func (node *ServiceSSE) Init(storage flownodes.Storage, form flownodes.Data) (err error) {
	if node.service, err = storage.GetService(node.ID); err != nil {
		return err
	}

	if node.Flow != "" {
		if node.flow, err = storage.GetFlow(node.Flow); err != nil {
			return err
		}
	}

	return
}

// Run send messages to the remote service and wait SSE stream response
// each sse message is sent to the to the selected flow
func (node *ServiceSSE) Run(c flownodes.Ctx) (resp flownodes.Data, err error) {
	if node.service == nil {
		return nil, errors.New("no service")
	}

	// build URL address
	url, err := node.service.GetAddress(node.Path)
	if err != nil {
		return resp, err
	}

	// create new http request
	req := flowdata.NewRequest()
	req.Method = "POST"
	req.URL = url
	req.Stream = true

	// set datasource id, flow id, flow run id as http headers
	req.Headers = c.Headers()

	// set body request
	req.Body = map[string]any{
		"model":    c.MustGetData("model"),
		"messages": c.MustGetData("messages"),
		"stream":   true,
	}

	// do request and receive response
	r, err := req.Do()
	if err != nil {
		return resp, err
	}
	defer r.Response.Body.Close()

	var answer string

	s := bufio.NewScanner(r.Response.Body)
	for s.Scan() {
		var msg sseMessage
		if err := json.Unmarshal(s.Bytes(), &msg); err != nil {
			return resp, err
		}

		// just for debug
		log.Printf("%+v", msg)

		data := map[string]any{"message": msg.Message}

		if s, ok := msg.Message["content"].(string); ok {
			answer += s
		}

		// trigger output ports
		c.Return(data)

		// flow Run submit data to the flow and wait response
		if node.flow != nil {
			_, err := node.flow.Run(data)
			if err != nil {
				return nil, err
			}
		}
	}

	// returns final flag when stream is over
	return flownodes.Data{"answer": answer}, nil
}

type sseMessage struct {
	Model     string         `json:"model"`
	CreatedAt time.Time      `json:"created_at"`
	Message   map[string]any `json:"message"`
	Done      bool           `json:"done"`

	// tail
	DoneReason         string        `json:"done_reason,omitempty"`
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuratuon time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}

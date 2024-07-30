package flownodes

import "github.com/osmoview/flowdata"

const NodeService = "service"

func init() {
	RegisterNode(NodeService, &Service{})
	SourceNodes = append(SourceNodes, NodeService)
}

// Service is a universal node for interacting with remote services
type Service struct {
	ID  string `yaml:"id"`
	srv RemoteService

	Input  []Port `yaml:"-"`
	Output []Port `yaml:"-"`
}

func (node *Service) Init(db DB) (err error) {
	if node.ID != "" {
		node.srv, err = db.GetService(node.ID)
		if err != nil {
			return err
		}

		node.Input, node.Output = node.srv.GetPorts()
	}

	return
}

func (node *Service) Execute(c Ctx) (out Resp, err error) {
	data := make(map[string]any)
	for _, p := range node.Input {
		data[p.Name] = c.MustGetData(p.Name)
	}

	resp := make(map[string]any)
	if node.srv != nil {
		if err := node.srv.Send(data, &resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (node *Service) GetPorts() (in, out []Port) {
	return node.Input, node.Output
}

func (node *Service) GetForm() flowdata.Form {
	return flowdata.Form{{
		Name:    "id",
		Label:   flowdata.L10n{"en": "Service", "ru": "Сервис"},
		Type:    flowdata.FieldTypeSelect,
		Options: []flowdata.FormFieldOption{},
		Hook:    flowdata.FieldHookServices,
	}}
}

func (node *Service) GetGroup() string {
	return GroupServices
}

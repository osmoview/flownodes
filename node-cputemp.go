package flownodes

import (
	"strconv"
	"strings"

	"github.com/ssimunic/gosensors"
)

type CPUtemp struct {
	Chip string `yaml:"chip"`
	Key  string `yaml:"key"`

	port Port
}

func (t *CPUtemp) Init(db DB) error {
	t.port = Port{Name: "temp", Type: PortNumber}
	return nil
}

func (t *CPUtemp) Execute(c Ctx) (out Resp, err error) {
	sensors, err := gosensors.NewFromSystem()
	if err != nil {
		return nil, err
	}

	var v string

	for name, chip := range sensors.Chips {
		if !strings.Contains(name, t.Chip) {
			continue
		}

		for key, val := range chip {
			if !strings.Contains(key, t.Key) {
				continue
			}

			v = val
			break
		}
	}

	v = strings.TrimPrefix(strings.TrimSuffix(v, "°C"), "+")
	temp, _ := strconv.ParseFloat(v, 64)

	return Resp{t.port.Name: temp}, nil
}

func (t *CPUtemp) GetPorts() (in, out []Port) {
	return nil, []Port{t.port}
}

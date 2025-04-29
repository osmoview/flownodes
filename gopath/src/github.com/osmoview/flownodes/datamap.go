package flownodes

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Data is just a `map[string]any`, for example, it used as response from the node
type Data map[string]any

// GetData from the map by the key
func (m Data) GetData(key string) (v any, ok bool) {
	if v, ok = m[key]; ok {
		return v, ok
	}

	return m.GetCaseInsensitive(key)
}

// MustGetData same as GetData, but do not returns second boolean variable
func (m Data) MustGetData(key string) any {
	if v, ok := m[key]; ok {
		return v
	}

	v, _ := m.GetCaseInsensitive(key)
	return v
}

func (m Data) GetCaseInsensitive(key string) (v any, ok bool) {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}

	return nil, false
}

// BindTo marshal-unmarshal data map to the struct by the yaml,
// yaml used, because node data stored in the yaml format
func (m Data) BindTo(x any) error {
	// no data no bind
	if len(m) == 0 {
		return nil
	}

	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, x)
}

func (m *Data) Set(k string, v any) {
	(*m)[k] = v
}

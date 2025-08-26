package flownodes

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
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

	cfg := &mapstructure.DecoderConfig{
		DecodeHook: func(f, t reflect.Type, data any) (any, error) {
			if s, ok := data.(string); ok && s == "" {
				return nil, nil
			}
			return data, nil
		},
		WeaklyTypedInput: true,
		Result:           x,
		TagName:          "yaml",
		Squash:           true,
		MatchName: func(mapKey, fieldName string) bool {
			if strings.EqualFold(mapKey, fieldName) {
				return true
			}

			mapKey = strings.ReplaceAll(mapKey, "_", "")
			return strings.EqualFold(mapKey, fieldName)
		},
	}

	dec, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}

	return dec.Decode(m.FixDotKeys())
}

func (m *Data) FixDotKeys() Data {
	data := Data{}

	for k, v := range *m {
		data.Set(k, v)
	}

	return data
}

func (m *Data) setmap(k string, v any) error {
	keyParts := strings.Split(k, ".")
	for i := 0; i < len(keyParts); i++ {
		key := keyParts[i]

		if i == len(keyParts)-1 {
			m.Set(key, v)
			break
		}

		sub, ok := (*m)[key]
		if !ok {
			sub = map[string]any{}
			(*m)[key] = sub
		}
		switch x := sub.(type) {
		case map[string]any:
			d := Data(x)
			return d.setmap(strings.Join(keyParts[i+1:], "."), v)
		default:
			return fmt.Errorf("%T is unexpected type", x)
		}
	}

	return nil
}

func (m *Data) Set(k string, v any) {
	if strings.Contains(k, ".") {
		m.setmap(k, v)
		return
	}

	(*m)[k] = v
}

func (m *Data) Has(k string) bool {
	_, ok := (*m)[k]
	return ok
}

func (m *Data) Contains(k string, v any) bool {
	vv, ok := (*m)[k]
	if !ok {
		return false
	}

	return vv == v
}

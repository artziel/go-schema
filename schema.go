package schema

import (
	"encoding/json"
)

type fieldErrors []error

func (se fieldErrors) MarshalJSON() ([]byte, error) {
	data := []byte("[")
	for i, err := range se {
		if i != 0 {
			data = append(data, ',')
		}

		j, err := json.Marshal(err.Error())
		if err != nil {
			return nil, err
		}

		data = append(data, j...)
	}
	data = append(data, ']')

	return data, nil
}

type Field struct {
	Errors fieldErrors `json:"errors"`
}

type Result struct {
	Fields map[string]Field `json:"fields"`
}

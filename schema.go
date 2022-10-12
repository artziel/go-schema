package schema

import (
	"encoding/json"
)

/*
Custom errors slice type. Used for extend json encoding field errors
*/
type fieldErrors []error

/*
Allow errors slice json encoding
*/
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

/*
Structure that will contain validation errors
*/
type Field struct {
	Errors fieldErrors `json:"errors"`
}

/*
Structure return by Validation function

After validation will contain only fields that fails
*/
type Result struct {
	Fields map[string]fieldErrors `json:"fields"`
}

/*
Check if Result structure contain at least one field with errors
*/
func (r *Result) HasErrors() bool {
	return len(r.Fields) > 0
}

/*
Add an error to a field on the Result structure. If fieldName do not exist on Result fields,
add a new field entry with the error, otherwise append the error to the field entry
*/
func (r *Result) AddFieldError(fieldName string, err error) {
	if err != nil {
		f := []error{}
		if _, exists := r.Fields[fieldName]; exists {
			f = append(r.Fields[fieldName], err)
		} else {
			f = []error{err}
		}
		r.Fields[fieldName] = f
	}
}

/*
Add an slice of errors to a field on the Result structure. If fieldName do not exist on Result fields,
add a new field entry with the error, otherwise append the errors to the field entry
*/
func (r *Result) AddFieldErrors(fieldName string, errs []error) {
	if errs != nil {
		if len(errs) > 0 {
			f := []error{}
			if _, exists := r.Fields[fieldName]; exists {
				f = r.Fields[fieldName]
				f = append(f, errs...)
			} else {
				f = errs
			}
			r.Fields[fieldName] = f
		}
	}
}

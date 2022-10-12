package schema

import (
	"encoding/json"
	"testing"
)

func TestSchema(t *testing.T) {
	result := Result{
		Fields: map[string]fieldErrors{},
	}

	if result.HasErrors() {
		t.Errorf("Test Fail, Result should not have errors")
	}

	result.AddFieldError("ID", ErrNumericFieldMaxValue)
	if !result.HasErrors() {
		t.Errorf("Test Fail, Result should have errors")
	}

	result.AddFieldError("ID", ErrNumericFieldMinValue)
	if !result.HasErrors() {
		t.Errorf("Test Fail, Result should have errors")
	}
	if len(result.Fields["ID"]) != 2 {
		if !result.HasErrors() {
			t.Errorf("Test Fail, Result should have 2 \"ID\" errors")
		}
	}

	result = Result{
		Fields: map[string]fieldErrors{},
	}
	errs := []error{ErrStringFieldMaxlength, ErrStringFieldrequired, ErrTagNameMissingValue}
	result.AddFieldErrors("Username", errs)
	if !result.HasErrors() || len(result.Fields["Username"]) != 3 {
		t.Errorf("Test Fail, Result should have 3 \"Username\" errors")
	}

	errs = []error{ErrTagRestrictToNotMatch, ErrTagMaxMissingValue}
	result.AddFieldErrors("Username", errs)
	if len(result.Fields["Username"]) != 5 {
		t.Errorf("Test Fail, Result should have 5 \"Username\" errors")
	}

	result = Result{
		Fields: map[string]fieldErrors{},
	}
	result.AddFieldError("ID", ErrNumericFieldMaxValue)
	encoded, _ := json.Marshal(result)
	expected := `{"fields":{"ID":["the field value is greater than the maximum numeric value allowed"]}}`
	if string(encoded) != expected {
		t.Errorf("Test Fail, Result JSON do not match expected value:\nGot  %s\nWant %s", string(encoded), expected)
	}
}

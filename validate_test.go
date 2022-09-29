package schema

import (
	"reflect"
	"testing"
)

func TestValidateTypes(t *testing.T) {
	type Model struct {
		ID        uint    `schema:"required"`
		Quantity  int     `schema:"min:1,max:5"`
		Size      float64 `schema:"require:Quantity Name"`
		Price     float64 `schema:"require:Quantity Name Error"`
		Shipping  float64 `schema:"min:12.21"`
		Margin    float64 `schema:"restrictTo:10.1 11.8 12.0"`
		Year      int     `schema:"require:XX Quantity Name"`
		Month     int     `schema:"restrictTo:1 2 3 4 5 6 7 8 9 10 11 12"`
		Day       int     `schema:"require:XX Quantity Name"`
		Name      string  `schema:"name:name,maxLength:5,minLength:1"`
		LastName  string  `schema:"required"`
		Cellphone string  `schema:"minLength:10"`
		SentAt    string  `schema:"require:ID"`
		Status    string  `schema:"restrictTo:ACTIVE SUSPENDED 'OUT OF STOCK'"`
		Delivered string  `schema:"restrictTo:YES NO"`
		Comment   string  `schema:"regex:[a-zA-Z].*"`
		BadRegex  string  `schema:"regex:\\"`
	}

	type Test struct {
		Input          Model
		ExpectedErrors map[string]map[error]struct{}
	}

	tests := []Test{
		{
			Input: Model{
				Quantity:  12,
				Name:      "This is a Test",
				Status:    "ACTIVE",
				Cellphone: "12345",
				Shipping:  2.4,
				Comment:   "123",
				BadRegex:  "123",
				Delivered: "NA",
				Price:     123.321,
				Year:      2022,
				Month:     12,
				SentAt:    "Home",
			},
			ExpectedErrors: map[string]map[error]struct{}{
				"ID":        {ErrNumericFieldrequired: struct{}{}},
				"Quantity":  {ErrNumericFieldMaxValue: struct{}{}},
				"Name":      {ErrStringFieldMaxlength: struct{}{}},
				"LastName":  {ErrStringFieldrequired: struct{}{}},
				"Cellphone": {ErrStringFieldMinlength: struct{}{}},
				"Comment":   {ErrStringFieldRegexMatchFail: struct{}{}},
				"BadRegex":  {ErrTagRegexValueFailToCompile: struct{}{}},
				"Price":     {ErrTagRequireFieldNotExists: struct{}{}},
				"Delivered": {ErrTagRestrictToNotMatch: struct{}{}},
				"Shipping":  {ErrNumericFieldMinValue: struct{}{}},
				"Year":      {ErrTagRequireFieldNotExists: struct{}{}},
				"Margin":    {ErrTagRestrictToNotMatch: struct{}{}},
				"SentAt":    {ErrTagRequireFieldFail: struct{}{}},
			},
		},
	}

	for i, test := range tests {
		// t Tag, model reflect.Value, field reflect.Value
		v := reflect.ValueOf(&test.Input)
		for j := 0; j < v.Elem().NumField(); j++ {
			errs := []error{}
			fieldName := v.Elem().Type().Field(j).Name
			tag := ParseTag(v.Elem().Type().Field(j).Tag.Get("schema"))

			switch v.Elem().Type().Field(j).Type.Kind() {
			case reflect.String:
				errs = validateString(tag, v, v.Elem().Field(j))
			case reflect.Float32, reflect.Float64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

				errs = validateNumeric(tag, v, v.Elem().Field(j))
			}

			if len(errs) != len(test.ExpectedErrors[fieldName]) {
				t.Errorf("Test %d unexpected result for field \"%s\":\nGot  %d Err\nWant %d Err", i, fieldName, len(errs), len(test.ExpectedErrors[fieldName]))
			} else {
				for _, e := range errs {
					if _, found := test.ExpectedErrors[fieldName][e]; !found {
						t.Errorf("Test %d unexpected error for field \"%s\": %s", i, fieldName, e)
					}
				}
			}
		}
	}

}

func TestValidate(t *testing.T) {
	type Model struct {
		ID        uint    `schema:"required"`
		Quantity  int     `schema:"min:1,max:5"`
		Size      float64 `schema:"require:Quantity Name"`
		Price     float64 `schema:"require:Quantity Name"`
		Shipping  float64 `schema:"min:12.21"`
		Margin    float64 `schema:"restrictTo:10.1 11.8 12.0"`
		Year      int     `schema:"require:Quantity Name"`
		Month     int     `schema:"restrictTo:1 2 3 4 5 6 7 8 9 10 11 12"`
		Day       int     `schema:"require:Quantity Name"`
		Name      string  `schema:"maxLength:5,minLength:1"`
		LastName  string  `schema:"required"`
		Cellphone string  `schema:"minLength:10"`
		SentAt    string  `schema:"require:ID"`
		Status    string  `schema:"restrictTo:ACTIVE SUSPENDED 'OUT OF STOCK'"`
		Delivered string  `schema:"restrictTo:YES NO"`
		Comment   string  `schema:"regex:[a-zA-Z].*"`
	}

	type Test struct {
		Input         Model
		ExpectedError error
	}

	tests := []Test{
		{
			Input: Model{
				Quantity:  2,
				Size:      12,
				Price:     128.4,
				Shipping:  14.0,
				Margin:    11.8,
				Year:      2022,
				Month:     12,
				Day:       24,
				Name:      "ASD",
				LastName:  "FDSA",
				Cellphone: "5555555555",
				SentAt:    "Home",
				Status:    "ACTIVE",
				Delivered: "YES",
				Comment:   "asdf",
			},
			ExpectedError: ErrValidationFail,
		},
		{
			Input: Model{
				ID:        1,
				Quantity:  2,
				Size:      12,
				Price:     128.4,
				Shipping:  14.0,
				Margin:    11.8,
				Year:      2022,
				Month:     12,
				Day:       24,
				Name:      "ASD",
				LastName:  "FDSA",
				Cellphone: "5555555555",
				SentAt:    "Home",
				Status:    "ACTIVE",
				Delivered: "YES",
				Comment:   "asdf",
			},
			ExpectedError: nil,
		},
	}

	if _, err := Validate(Model{}); err == nil {
		t.Errorf("Test expected error validation:\nGot  No Error\nWant %s", ErrValidatePtrExpected)
	}

	for i, test := range tests {
		result, err := Validate(&test.Input)
		if err == nil && test.ExpectedError != nil {
			t.Errorf("Test %d expected error validation:\nGot  No Error\nWant %s", i, test.ExpectedError)
		} else if err != nil && test.ExpectedError == nil {
			t.Errorf("Test %d unexpected error validation:\nGot  %s\nWant No Error", i, err)
		} else if err != nil && test.ExpectedError != err {
			t.Errorf("Test %d unexpected error validation:\nGot  %s\nWant %s", i, err, test.ExpectedError)
		} else if err != nil && test.ExpectedError == err {
			if !result.HasErrors() {
				t.Errorf("Test %d validation result should have errors", i)
			}
		}
	}

}

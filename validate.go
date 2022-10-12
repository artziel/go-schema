package schema

import (
	"reflect"
	"regexp"
)

/*
Validate struct string field

Rules:

	    required:	Empty value not allow
	    minLength:	Min string length allow
	    maxLength:	Max string length allow
	    regex:	Value must match with regular expresion
		restrictTo: Field value must be one of the listed values

An slice of type "error" will be return containing each failed rule
*/
func validateString(t Tag, model reflect.Value, field reflect.Value) []error {
	errs := []error{}

	value := field.String()

	if t.Exists("required") && value == "" {
		errs = append(errs, ErrStringFieldrequired)
	} else if value != "" {
		length := len([]rune(value))
		if t.GetUint("minLength") > 0 && uint64(length) < t.GetUint("minLength") {
			errs = append(errs, ErrStringFieldMinlength)
		}
		if t.GetUint("maxLength") > 0 && uint64(length) > t.GetUint("maxLength") {
			errs = append(errs, ErrStringFieldMaxlength)
		}
		if t.HasValue("regex") {
			r, err := regexp.Compile(t.GetString("regex"))
			if err != nil {
				errs = append(errs, ErrTagRegexValueFailToCompile)
			} else if !r.MatchString(value) {
				errs = append(errs, ErrStringFieldRegexMatchFail)
			}

		}
	}

	values := t.GetSliceString("restrictTo", ' ')
	if len(values) > 0 {
		match := false
		for _, v := range values {
			if v == value {
				match = true
			}
		}
		if !match {
			errs = append(errs, ErrTagRestrictToNotMatch)
		}
	}

	if field.Interface() != reflect.Zero(field.Type()).Interface() {
		if err := validateRequirements(model, t); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

/*
Validate struct numeric field

Rules:

	    required:	Empty value not allow
	    min:	Min numeric value allow
	    max:	Max numeric value allow
		restrictTo: Field value must be one of the listed values

An slice of type "error" will be return containing each failed rule
*/
func validateNumeric(t Tag, model reflect.Value, field reflect.Value) []error {
	errs := []error{}

	var value float64

	switch field.Kind() {
	case reflect.Float32, reflect.Float64:
		value = field.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = float64(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = float64(field.Uint())
	}

	if t.Exists("required") && value == 0 {
		errs = append(errs, ErrNumericFieldrequired)
	} else if value != 0 {
		if t.GetFloat("min") > 0 && value < t.GetFloat("min") {
			errs = append(errs, ErrNumericFieldMinValue)
		}
		if t.GetFloat("max") > 0 && value > t.GetFloat("max") {
			errs = append(errs, ErrNumericFieldMaxValue)
		}
	}

	values := t.GetSliceFloat("restrictTo", ' ')
	if len(values) > 0 {
		match := false
		for _, v := range values {
			if v == value {
				match = true
			}
		}
		if !match {
			errs = append(errs, ErrTagRestrictToNotMatch)
		}
	}

	if field.Interface() != reflect.Zero(field.Type()).Interface() {
		if err := validateRequirements(model, t); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

/*
Validate if the fields required by a field are empty and exist in the structure
*/
func validateRequirements(model reflect.Value, t Tag) error {
	var err error
	if t.Exists("require") {
		for _, r := range t.GetSliceString("require", ' ') {
			_, exists := model.Elem().Type().FieldByName(r)
			if !exists {
				err = ErrTagRequireFieldNotExists
				break
			} else {
				required := reflect.Indirect(model).FieldByName(r)
				if required.Interface() == reflect.Zero(required.Type()).Interface() {
					err = ErrTagRequireFieldFail
					break
				}
			}
		}
	}
	return err
}

/*
Validate each structure field with a valid schema tag. The parameter model must
be an struct ptr
*/
func Validate(model interface{}) (Result, error) {
	result := Result{
		Fields: map[string]fieldErrors{},
	}

	v := reflect.ValueOf(model)

	if v.Kind() != reflect.Ptr {
		return result, ErrValidatePtrExpected
	}

	for i := 0; i < v.Elem().NumField(); i++ {
		t := ParseTag(v.Elem().Type().Field(i).Tag.Get("schema"))

		if !t.IsEmpty() {

			fieldName := t.GetString("name")
			if fieldName == "" {
				fieldName = v.Elem().Type().Field(i).Name
			}

			switch v.Elem().Type().Field(i).Type.Kind() {
			case reflect.String:
				result.AddFieldErrors(fieldName, validateString(t, v, v.Elem().Field(i)))
			case reflect.Float32, reflect.Float64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

				result.AddFieldErrors(fieldName, validateNumeric(t, v, v.Elem().Field(i)))
			}
		}
	}

	if result.HasErrors() {
		return result, ErrValidationFail
	}

	return result, nil
}

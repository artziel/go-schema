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
func validateString(t Tag, value string) []error {
	errs := []error{}

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
		if t.GetString("regex") != "" {
			r, err := regexp.Compile(value)
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
func validateNumeric(t Tag, value float64) []error {
	errs := []error{}

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

	require := t.GetSliceFloat("restrictTo", ' ')
	if len(require) > 0 {
		match := false
		for _, v := range require {
			if v == value {
				match = true
			}
		}
		if !match {
			errs = append(errs, ErrTagRestrictToNotMatch)
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
		Fields: map[string]Field{},
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
				value := v.Elem().Field(i).String()
				result.AddFieldErrors(fieldName, validateString(t, value))
			case reflect.Float32, reflect.Float64:
				value := v.Elem().Field(i).Float()
				result.AddFieldErrors(fieldName, validateNumeric(t, value))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value := v.Elem().Field(i).Int()
				result.AddFieldErrors(fieldName, validateNumeric(t, float64(value)))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value := v.Elem().Field(i).Uint()
				result.AddFieldErrors(fieldName, validateNumeric(t, float64(value)))
			}
			if v.Elem().Field(i).Interface() != reflect.Zero(v.Elem().Field(i).Type()).Interface() {
				if err := validateRequirements(v, t); err != nil {
					result.AddFieldError(fieldName, err)
				}
			}
		}
	}

	if result.HasErrors() {
		return result, ErrValidationFail
	}

	return result, nil
}

package schema

import (
	"fmt"
	"reflect"
	"regexp"
)

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

	require := t.GetSliceString("restrictTo", ' ')
	fmt.Printf("%v\n", require)
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

			fieldErrors := []error{}
			checkRequirements := false

			switch v.Elem().Type().Field(i).Type.Kind() {
			case reflect.String:
				value := v.Elem().Field(i).String()
				fieldErrors = append(fieldErrors, validateString(t, value)...)
				if value != "" {
					checkRequirements = true
				}
			case reflect.Float32, reflect.Float64:
				value := v.Elem().Field(i).Float()
				fieldErrors = append(fieldErrors, validateNumeric(t, value)...)
				if value != 0 {
					checkRequirements = true
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value := v.Elem().Field(i).Int()
				fieldErrors = append(fieldErrors, validateNumeric(t, float64(value))...)
				if value != 0 {
					checkRequirements = true
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value := v.Elem().Field(i).Uint()
				fieldErrors = append(fieldErrors, validateNumeric(t, float64(value))...)
				if value != 0 {
					checkRequirements = true
				}
			}

			if t.Exists("require") && checkRequirements {
				for _, r := range t.GetSliceString("require", ' ') {
					_, exists := v.Elem().Type().FieldByName(r)
					if !exists {
						fieldErrors = append(fieldErrors, ErrTagRequireFieldNotExists)
						break
					} else {
						required := reflect.Indirect(v).FieldByName(r)
						if required.Kind() == reflect.String {
							if required.String() == "" {
								fieldErrors = append(fieldErrors, ErrTagRequireFieldFail)
								break
							}
						} else if required.Kind() == reflect.Float32 || required.Kind() == reflect.Float64 {
							if required.Float() == 0 {
								fieldErrors = append(fieldErrors, ErrTagRequireFieldFail)
								break
							}
						} else if required.Kind() == reflect.Int || required.Kind() == reflect.Int8 || required.Kind() == reflect.Int16 || required.Kind() == reflect.Int32 || required.Kind() == reflect.Int64 {
							if required.Int() == 0 {
								fieldErrors = append(fieldErrors, ErrTagRequireFieldFail)
								break
							}
						} else if required.Kind() == reflect.Uint || required.Kind() == reflect.Uint8 || required.Kind() == reflect.Uint16 || required.Kind() == reflect.Uint32 || required.Kind() == reflect.Uint64 {
							if required.Uint() == 0 {
								fieldErrors = append(fieldErrors, ErrTagRequireFieldFail)
								break
							}
						}
					}
				}
			}

			if len(fieldErrors) > 0 {
				if t.GetString("name") != "" {
					result.Fields[t.GetString("name")] = Field{Errors: fieldErrors}
				} else {
					result.Fields[v.Elem().Type().Field(i).Name] = Field{Errors: fieldErrors}
				}
			}
		}
	}

	if len(result.Fields) > 0 {
		return result, ErrValidationFail
	}

	return result, nil
}

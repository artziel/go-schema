package schema

import "reflect"

func Validate(model interface{}) (Result, error) {
	result := Result{
		Fields: map[string]Field{},
	}

	v := reflect.ValueOf(model)

	if v.Kind() != reflect.Ptr {
		return result, ErrValidatePtrExpected
	}

	for i := 0; i < v.Elem().NumField(); i++ {
		tag := v.Elem().Type().Field(i).Tag.Get("schema")
		if tag != "" {

			field := v.Elem().Field(i)
			fieldName := v.Elem().Type().Field(i).Name
			fieldErrors := []error{}

			st, err := parse_tag(tag)
			if err != nil {
				fieldErrors = append(fieldErrors, err)
			} else {
				kind := v.Elem().Type().Field(i).Type.Kind()
				switch kind {
				case reflect.String:
					value := field.String()
					if st.Required && value == "" {
						fieldErrors = append(fieldErrors, ErrStringFieldrequired)
					} else {
						length := len([]rune(value))
						if st.MinLength > 0 && uint(length) < st.MinLength {
							fieldErrors = append(fieldErrors, ErrStringFieldMinlength)
						}
						if st.MaxLength > 0 && uint(length) > st.MaxLength {
							fieldErrors = append(fieldErrors, ErrStringFieldMaxlength)
						}
						if st.Regex != nil && !st.Regex.MatchString(value) {
							fieldErrors = append(fieldErrors, ErrStringFieldRegexMatchFail)
						}
					}
				case reflect.Float32, reflect.Float64:
					value := field.Float()
					if st.Required && value == 0 {
						fieldErrors = append(fieldErrors, ErrNumericFieldrequired)
					} else {
						if st.Min > 0 && value < st.Min {
							fieldErrors = append(fieldErrors, ErrNumericFieldMinValue)
						}
						if st.Max > 0 && value > st.Max {
							fieldErrors = append(fieldErrors, ErrNumericFieldMaxValue)
						}
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value := field.Int()
					if st.Required && value == 0 {
						fieldErrors = append(fieldErrors, ErrNumericFieldrequired)
					} else {
						if st.Min > 0 && float64(value) < st.Min {
							fieldErrors = append(fieldErrors, ErrNumericFieldMinValue)
						}
						if st.Max > 0 && float64(value) > st.Max {
							fieldErrors = append(fieldErrors, ErrNumericFieldMaxValue)
						}
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value := field.Uint()
					if st.Required && value == 0 {
						fieldErrors = append(fieldErrors, ErrNumericFieldrequired)
					} else {
						if st.Min > 0 && float64(value) < st.Min {
							fieldErrors = append(fieldErrors, ErrNumericFieldMinValue)
						}
						if st.Max > 0 && float64(value) > st.Max {
							fieldErrors = append(fieldErrors, ErrNumericFieldMaxValue)
						}
					}
				}

				if len(st.Require) > 0 {
					for _, r := range st.Require {
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
			}
			if len(fieldErrors) > 0 {
				if st.Name != "" {
					result.Fields[st.Name] = Field{Errors: fieldErrors}
				} else {
					result.Fields[fieldName] = Field{Errors: fieldErrors}
				}
			}
		}
	}

	if len(result.Fields) > 0 {
		return result, ErrValidationFail
	}

	return result, nil
}

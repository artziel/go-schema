package schema

import (
	"errors"
	"fmt"
)

type SchemaError struct {
	Code int
	Err  error
}

func (e *SchemaError) Error() string {
	return fmt.Sprintf("(%d) %v", e.Code, e.Err)
}

func NewError(code int, message string) *SchemaError {
	return &SchemaError{
		Code: code,
		Err:  errors.New(message),
	}
}

var (
	ErrValidatePtrExpected = NewError(1001, "function \"schema.Validate\" expect a pointer")
	ErrValidationFail      = NewError(1002, "schema validation fail")
)

var (
	ErrTagRequireMissingValue   = NewError(2001, "missing value for \"require\" parameter")
	ErrTagRequireFieldNotExists = NewError(2002, "one or more of the fields required by the field do not exist in the structure")
	ErrTagRequireFieldFail      = NewError(2003, "one or more of the fields required are empty")
)

var (
	ErrTagRestrictToMissingValue = NewError(3001, "missing value for \"restrictTo\" parameter")
	ErrTagRestrictToNotMatch     = NewError(3002, "value do not match any \"restrictTo\" parameter values")
)

var (
	ErrTagRegexMissingValue       = NewError(4001, "missing value for \"regex\" parameter")
	ErrTagRegexValueFailToCompile = NewError(4002, "value for \"regex\" parameter fail to compile")
)

var (
	ErrTagNameMissingValue = NewError(5001, "missing value for \"name\" parameter")
	ErrTagNameValueInvalid = NewError(5002, "value for \"name\" is invalid")
)

var (
	ErrTagMinLengthMissingValue = NewError(6001, "missing value for \"minLength\" parameter")
	ErrTagMaxLengthMissingValue = NewError(6002, "missing value for \"maxLength\" parameter")
	ErrTagMinMissingValue       = NewError(6003, "missing value for \"min\" parameter")
	ErrTagMaxMissingValue       = NewError(6004, "missing value for \"max\" parameter")
)

var (
	ErrStringFieldRegexMatchFail = NewError(7001, "field value failed regular expression evaluation")
	ErrStringFieldrequired       = NewError(7002, "field required, not empty string allow")
	ErrStringFieldMinlength      = NewError(7003, "the length of the field is less than the minimum length allowed")
	ErrStringFieldMaxlength      = NewError(7004, "the length of the field is greater than the maximum length allowed")
)

var (
	ErrNumericFieldrequired = NewError(8001, "field required, 0 value not allow")
	ErrNumericFieldMinValue = NewError(8002, "field value is less than the minimum numeric value allowed")
	ErrNumericFieldMaxValue = NewError(8003, "the field value is greater than the maximum numeric value allowed")
)

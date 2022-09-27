package schema

import "errors"

var ErrValidatePtrExpected = errors.New("function \"shcema.Validate\" expect a pointer")
var ErrValidationFail = errors.New("schema validation fail")

var ErrTagRequireMissingValue = errors.New("missing value for \"require\" parameter")
var ErrTagRequireFieldNotExists = errors.New("one or more of the fields required by the field do not exist in the structure")
var ErrTagRequireFieldFail = errors.New("one or more of the fields required by the field are empty")

var ErrTagRestrictToMissingValue = errors.New("missing value for \"restrictTo\" parameter")
var ErrTagRestrictToNotMatch = errors.New("value do not match any \"restrictTo\" parameter values")

var ErrTagRegexMissingValue = errors.New("missing value for \"regex\" parameter")
var ErrTagRegexValueFailToCompile = errors.New("value for \"regex\" parameter fail to compile")

var ErrTagNameMissingValue = errors.New("missing value for \"name\" parameter")
var ErrTagNameValueInvalid = errors.New("value for \"name\" is invalid")

var ErrTagMinLengthMissingValue = errors.New("missing value for \"minLength\" parameter")
var ErrTagMaxLengthMissingValue = errors.New("missing value for \"maxLength\" parameter")
var ErrTagMinMissingValue = errors.New("missing value for \"min\" parameter")
var ErrTagMaxMissingValue = errors.New("missing value for \"max\" parameter")

var ErrStringFieldRegexMatchFail = errors.New("field value failed regular expression evaluation")
var ErrStringFieldrequired = errors.New("field required, not empty string allow")
var ErrStringFieldMinlength = errors.New("the length of the field is less than the minimum length allowed")
var ErrStringFieldMaxlength = errors.New("the length of the field is greater than the maximum length allowed")

var ErrNumericFieldrequired = errors.New("field required, 0 value not allow")
var ErrNumericFieldMinValue = errors.New("field value is less than the minimum numeric value allowed")
var ErrNumericFieldMaxValue = errors.New("the field value is greater than the maximum numeric value allowed")

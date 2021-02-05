package taskdef

import (
	"overseer/common/validator"
	"regexp"

	vl "github.com/go-playground/validator/v10"
)

func init() {

	validator.Valid.RegisterValidatorRule("resname", ResourceNameValidator)
	validator.Valid.RegisterValidatorRule("varname", VariableNameValidator)

}

//ResourceNameValidator - validator function  for resource name like task, flag or ticket
func ResourceNameValidator(fl vl.FieldLevel) bool {

	if val, ok := fl.Field().Interface().(string); ok {
		result, _ := validateValueResource(val)
		return result
	}

	return false
}

//VariableNameValidator - validator function  for resource name like task, flag or ticket
func VariableNameValidator(fl vl.FieldLevel) bool {

	if val, ok := fl.Field().Interface().(string); ok {
		result, _ := validateValueVariable(val)
		return result
	}

	return false
}

func validateValueResource(resource string) (bool, error) {

	return regexp.MatchString(`^[A-Za-z][\w\-\.]*$`, resource)

}

func validateValueVariable(resource string) (bool, error) {

	return regexp.MatchString(`^%%[\dA-Z]*$`, resource)

}

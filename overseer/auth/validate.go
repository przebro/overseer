package auth

import (
	"overseer/common/validator"
	"regexp"

	vl "github.com/go-playground/validator/v10"
)

func init() {
	validator.Valid.RegisterValidatorRule("auth", NameAuthValidator)
}

//NameAuthValidator - validator function  for  description
func NameAuthValidator(fl vl.FieldLevel) bool {

	if val, ok := fl.Field().Interface().(string); ok {
		result, _ := validateValueAuth(val)
		return result
	}

	return false
}
func validateValueAuth(resource string) (bool, error) {

	return regexp.MatchString(`^\w[\w" "()\.]+$`, resource)

}

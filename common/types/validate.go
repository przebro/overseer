package types

import (
	"fmt"
	"goscheduler/common/validator"
	"regexp"

	vl "github.com/go-playground/validator/v10"
)

func init() {
	validator.Valid.RegisterTypeValidator("HourMinTime", "hmtime", HourMinTimeValidator)
}

//HourMinTimeValidator - validator function for HH:MM format
func HourMinTimeValidator(fl vl.FieldLevel) bool {

	xx := fl.Field().Interface()
	fmt.Println(xx)
	if val, ok := fl.Field().Interface().(HourMinTime); ok {
		result, _ := validateValueHourMinTime(val)
		return result
	}

	return false

}

func validateValueHourMinTime(resource HourMinTime) (bool, error) {

	return regexp.MatchString(`^(?:([01]?\d|2[0-3]):([0-5]?\d))$`, string(resource))

}

package unique

import (
	"overseer/common/validator"
	"regexp"

	vl "github.com/go-playground/validator/v10"
)

func init() {

	validator.Valid.RegisterTypeValidator("TaskOrderID", "orderID", TaskIDValidator)
}

//TaskIDValidator - Validator function for a TaskOrederID type
func TaskIDValidator(fl vl.FieldLevel) bool {

	if actual, ok := fl.Field().Interface().(TaskOrderID); ok {
		result, _ := actual.validateValue()
		return result
	}

	return false

}

//ValidateValue - Validates TaskOrderID
func (orderID TaskOrderID) validateValue() (bool, error) {

	if len(orderID) != 5 {
		return false, errInvalidLen
	}
	match, _ := regexp.MatchString(`[0-9A-Za-z]{5}`, string(orderID))

	if !match {
		return false, errInvalidChar
	}
	return true, nil
}

package unique

import (
	"overseer/common/validator"

	vl "github.com/go-playground/validator/v10"
)

func init() {

	validator.Valid.RegisterTypeValidator("TaskOrderID", "orderID", TaskIDValidator)
}

//TaskIDValidator - Validator function for a TaskOrederID type
func TaskIDValidator(fl vl.FieldLevel) bool {

	actual := fl.Field().Interface().(TaskOrderID)
	ok, _ := actual.validateValue()
	return ok

}

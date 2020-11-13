package validator

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type DataValidator struct {
	v        *validator.Validate
	regTypes map[string]string
}

var Valid *DataValidator
var errValidatorNotRegistered = errors.New("Validator for given type not registered")

func init() {
	Valid = &DataValidator{v: validator.New(), regTypes: make(map[string]string, 0)}
}

//RegisterValidatorRule - Registers a custom rule for a field validation
func (dv *DataValidator) RegisterValidatorRule(rule string, vfunc validator.Func) error {

	err := dv.v.RegisterValidation(rule, vfunc, false)
	return err

}

/*RegisterValidator - Registers a custom rule for type validation. This method should be used only for
types that can't be described with struct tag e.g. variables of custom type
*/
func (dv *DataValidator) RegisterTypeValidator(typeName, rule string, vfunc validator.Func) error {

	err := dv.RegisterValidatorRule(rule, vfunc)

	if err == nil {
		dv.regTypes[typeName] = rule
	}

	return err
}

func (dv *DataValidator) ValidateTag(s interface{}, tag string) error {

	return dv.v.Var(s, tag)
}

func (dv *DataValidator) Validate(s interface{}) error {
	var err error

	if reflect.TypeOf(s).Kind() == reflect.Struct {
		err = dv.v.Struct(s)
	} else {
		typeName := reflect.TypeOf(s).Name()
		v, exists := dv.regTypes[typeName]
		if !exists {
			return errValidatorNotRegistered
		}

		err = dv.v.Var(s, v)
	}

	return err
}

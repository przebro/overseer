package date

import (
	"goscheduler/common/validator"

	vl "github.com/go-playground/validator/v10"
)

func init() {

	validator.Valid.RegisterTypeValidator("Odate", "odate", OdateValidator)
	validator.Valid.RegisterTypeValidator("OdateValue", "odateval", OdateValueValidator)
}

//OdateValidator - validator function for a Odate type
func OdateValidator(fl vl.FieldLevel) bool {

	actual := fl.Field().Interface().(Odate)
	ok, _ := actual.validateValue()
	return ok

}

//OdateValueValidator - validator function for a OdateValue type
func OdateValueValidator(fl vl.FieldLevel) bool {

	actual := fl.Field().Interface().(OdateValue)
	ok, _ := actual.validateValue()
	return ok

}

//ValidateValue - Validates if given odate is a valid date
func (date Odate) validateValue() (bool, error) {

	day31 := map[int]int{1: 31, 3: 31, 5: 31, 7: 31, 8: 31, 10: 31, 12: 31}

	if date == OdateNone {
		return true, nil
	}

	if len(date) != 8 {
		return false, errOdateInvalidLen
	}

	for _, v := range date {

		if v < '0' || v > '9' {
			return false, errOdateNotNumeric
		}
	}

	y, m, d := date.Ymd()

	if y == 0 {
		return false, errOdateInvalidYear
	}

	leapYear := false

	if (y%4 == 0) && (y%100 != 0) || (y%400 == 0) {

		leapYear = true
	}

	if m < 1 || m > 12 {
		return false, errOdateInvalidMonth
	}

	if d < 1 || d > 31 {
		return false, errOdateInvalidDay
	}

	if m == 2 && ((leapYear && d > 29) || (!leapYear && d > 28)) {
		return false, errOdateInvalidDay
	}

	if v, is := day31[m]; is && d > v || !is && d > 30 {
		return false, errOdateInvalidDay
	}

	return true, nil
}

func (oval OdateValue) validateValue() (bool, error) {
	if oval == OdateValueNone || oval == OdateValueAny {
		return true, nil
	}

	if oval == OdateValueDate || oval == OdateValueNext || oval == OdateValuePrev {
		return true, nil
	}

	err := validator.Valid.Validate(Odate(oval))
	if err != nil {
		return false, err
	}

	return true, nil
}

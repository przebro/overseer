package date

import (
	"testing"

	"github.com/przebro/overseer/common/validator"
)

func TestOdateValidator(t *testing.T) {

	value := Odate("20200132")
	err := validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = Odate("20200112")

	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

}
func TestOdateValueValidator(t *testing.T) {

	value := OdateValue("")
	err := validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("*")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("ODATE")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("NEXT")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("PREV")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("20201105")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("?")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("20200230")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = OdateValue(" ")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("-001")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("+001")
	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("=001")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("+01")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("001")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = OdateValue("+00A")
	err = validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

}

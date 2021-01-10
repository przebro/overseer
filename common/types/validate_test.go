package types

import (
	"overseer/common/validator"
	"testing"
)

func TestHourMinTime(t *testing.T) {

	tm := HourMinTime("23:20")
	err := validator.Valid.Validate(tm)
	if err != nil {
		t.Error(err)
	}

	tm = HourMinTime("00:00")
	err = validator.Valid.Validate(tm)
	if err != nil {
		t.Error(err)
	}

	tm = HourMinTime("00:60")
	err = validator.Valid.Validate(tm)
	if err == nil {
		t.Error("unexpected value")
	}
	tm = HourMinTime("24:00")
	err = validator.Valid.Validate(tm)
	if err == nil {
		t.Error("unexpected value")
	}

}

package types

import (
	"testing"

	"github.com/przebro/overseer/common/validator"
)

type teststruct struct {
	Tm string `validate:"hmtime"`
}

func TestHourMinTime(t *testing.T) {

	tst := teststruct{}
	err := validator.Valid.Validate(tst)
	if err == nil {
		t.Error("unexpected value")
	}

	tm := HourMinTime("23:20")
	err = validator.Valid.Validate(tm)
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

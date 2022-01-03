package auth

import (
	"testing"

	"github.com/przebro/overseer/common/validator"
)

type teststruct struct {
	InvalidField int `validate:"auth"`
}

func TestValidator(t *testing.T) {

	tdata := teststruct{}

	err := validator.Valid.Validate(tdata)
	if err == nil {
		t.Error("unexpected result")
	}

	err = validator.Valid.ValidateTag("Abcdef0123_ .", "auth")
	if err != nil {
		t.Error("unexpected result")
	}

	err = validator.Valid.ValidateTag("Abcdef0123_ #", "auth")
	if err == nil {
		t.Error("unexpected result")
	}

	err = validator.Valid.ValidateTag("Abc%def0123_", "auth")
	if err == nil {
		t.Error("unexpected result")
	}

}

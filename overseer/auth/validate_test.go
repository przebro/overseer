package auth

import (
	"overseer/common/validator"
	"testing"
)

func TestValidator(t *testing.T) {
	err := validator.Valid.ValidateTag("Abcdef0123_ .", "auth")
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

package unique

import (
	"testing"

	"github.com/przebro/overseer/common/validator"
)

func TestValidator(t *testing.T) {

	value := TaskOrderID("123456")
	err := validator.Valid.Validate(value)
	if err == nil {
		t.Error("Unexpected value")
	}

	value = TaskOrderID("12345")

	err = validator.Valid.Validate(value)
	if err != nil {
		t.Error("Unexpected value")
	}

}

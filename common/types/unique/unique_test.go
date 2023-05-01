package unique

import (
	"encoding/hex"
	"testing"
)

func TestUnique(t *testing.T) {

	result := NewID()

	bytes := make([]byte, 12)
	copy(bytes, result[:])

	str := hex.EncodeToString(bytes)

	if str != result.Hex() {
		t.Error("hex encode error")
	}
	none := None()
	for i := range none {
		if none[i] != 0 {
			t.Error("none value not equal zero")
		}
	}
}

func TestValidateOrderID(t *testing.T) {

	ok, err := TaskOrderID("00000").validateValue()
	if !ok {
		t.Error(err)
	}

	ok, err = TaskOrderID("1111").validateValue()
	if ok {
		t.Error("Unexpected value, expected", false)
	}

	if err != errInvalidLen {
		t.Error("Unexpected value, expected:", errInvalidLen, "actual:", err)
	}

	ok, err = TaskOrderID("112211").validateValue()
	if ok {
		t.Error("Unexpected value, expected", false)
	}

	if err != errInvalidLen {
		t.Error("Unexpected value, expected:", errInvalidLen, "actual:", err)
	}

	ok, err = TaskOrderID("1aA2_").validateValue()
	if ok {
		t.Error("Unexpected value, expected", false)
	}

	if err != errInvalidChar {
		t.Error("Unexpected value, expected:", errInvalidChar, "actual:", err)
	}

}

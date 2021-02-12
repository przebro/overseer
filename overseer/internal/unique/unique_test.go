package unique

import (
	"encoding/hex"
	"overseer/common/types/date"
	"testing"
)

func TestSequence(t *testing.T) {

	//sequence of order id starts from 1

	week := date.CurrentOdate().Woyear()
	for x := 0; x < 3845; x++ {
		result := NewOrderID()
		//0,0,61 = 0,0,Z
		if x == 60 && string(result) != string([]byte{base62Str[week], '0', '0', '0', 'Z'}) {
			t.Error("invalid order ID, unexpected value", result)
		}
		// 0*62^2, 61*62^1,0 = 0,Z,0
		if x == 3781 && string(result) != string([]byte{base62Str[week], '0', '0', 'Z', '0'}) {
			t.Error("invalid order ID, unexpected value", result)
		}
		// 1* 62^2,0*62^1,0 = 1,0,0
		if x == 3845 && string(result) != string([]byte{base62Str[week], '0', '1', '0', '0'}) {
			t.Error("invalid order ID, unexpected value", result)
		}
	}

}

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

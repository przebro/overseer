package unique

import (
	"encoding/hex"
	"goscheduler/overseer/internal/date"
	"testing"
)

func TestSequence(t *testing.T) {

	//sequence of order id starts from 1

	month := date.CurrentOdate().Omonth()
	for x := 0; x < 3845; x++ {
		result := NewOrderID()
		//0,0,61 = 0,0,Z
		if x == 60 && string(result) != string([]byte{month[0], month[1], '0', '0', 'Z'}) {
			t.Error("invalid order ID, unexpected value", result)
		}
		// 0*62^2, 61*62^1,0 = 0,Z,0
		if x == 3781 && string(result) != string([]byte{month[0], month[1], '0', 'Z', '0'}) {
			t.Error("invalid order ID, unexpected value", result)
		}
		// 1* 62^2,0*62^1,0 = 1,0,0
		if x == 3845 && string(result) != string([]byte{month[0], month[1], '1', '0', '0'}) {
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

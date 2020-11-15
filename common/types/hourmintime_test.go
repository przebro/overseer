package types

import "testing"

func TestToString(t *testing.T) {

	if "10:30" != HourMinTime("10:30").String() {
		t.Fatal("fatal")
	}
}

func TestAsTime(t *testing.T) {

	h, m := HourMinTime("10:30").AsTime()

	if h != 10 || m != 30 {
		t.Error("Unexpected values", 10, 30)
	}

}

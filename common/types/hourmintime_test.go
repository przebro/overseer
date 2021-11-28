package types

import (
	"fmt"
	"testing"
	"time"
)

func TestToString(t *testing.T) {

	if HourMinTime("10:30").String() != "10:30" {
		t.Fatal("fatal")
	}
}

func TestAsTime(t *testing.T) {

	h, m := HourMinTime("10:30").AsTime()

	if h != 10 || m != 30 {
		t.Error("Unexpected values", 10, 30)
	}
}

func TestNow(t *testing.T) {

	r := Now()
	fmt.Println(r)
}

func TestFromTime(t *testing.T) {
	tm := time.Now()
	result := FromTime(tm)
	if fmt.Sprintf("%02d:%02d", tm.Hour(), tm.Minute()) != result.String() {
		t.Error("unexpected result, invalid date:", result.String())
	}
}

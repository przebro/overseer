package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//HourMinTime - represents formatted time value
type HourMinTime string

func (h HourMinTime) String() string {

	return string(h)
}

//AsTime - returns HourMinTimne as hour,min
func (h HourMinTime) AsTime() (hour, min int) {
	res := strings.Split(string(h), ":")
	hour, _ = strconv.Atoi(res[0])
	min, _ = strconv.Atoi(res[1])
	return
}

func Now() HourMinTime {

	hmt := time.Now()
	return HourMinTime(fmt.Sprintf("%02d:%02d", hmt.Hour(), hmt.Minute()))
}

func FromTime(tm time.Time) HourMinTime {

	return HourMinTime(fmt.Sprintf("%02d:%02d", tm.Hour(), tm.Minute()))
}

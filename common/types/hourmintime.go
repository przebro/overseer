package types

import (
	"strconv"
	"strings"
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

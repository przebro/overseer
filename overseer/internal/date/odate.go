package date

import (
	"fmt"
	"goscheduler/overseer/internal/taskdef"
	"strconv"
	"time"
)

//Odate - helper type for order date format YYYYMMDD
type Odate string

//ODATE -returns odate in format YYMMDD
func (date Odate) ODATE() string {

	return string(date[2:])
}

//Oday - returns day from odate
func (date Odate) Oday() string {
	return string(date[6:])
}

//Omonth - returns month from odate
func (date Odate) Omonth() string {
	return string(date[4:6])
}

//Oyear - returns last two digits of year from odate
func (date Odate) Oyear() string {
	return string(date[2:4])
}

//Ocent - returns century from odate
func (date Odate) Ocent() string {
	return string(date[0:2])
}

//FormatDate - formats odate into format YYYY-MM-DD
func (date Odate) FormatDate() string {
	return fmt.Sprintf("%s-%s-%s", string(date[0:4]), string(date[4:6]), string(date[6:]))
}

//Ymd - Returns int values of given Odate
func (date Odate) Ymd() (year, month, day int) {
	year, _ = strconv.Atoi(string(date[0:4]))
	month, _ = strconv.Atoi(string(date[4:6]))
	day, _ = strconv.Atoi(string(date[6:]))

	return
}

//Doyear - returns the day of year from odate
func (date Odate) Doyear() int {
	year, _ := strconv.Atoi(fmt.Sprintf("20%s", date.Oyear()))
	month, _ := strconv.Atoi(date.Omonth())
	day, _ := strconv.Atoi(date.Oday())

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return t.YearDay()
}

//Woyear - returns week number from odate
func (date Odate) Woyear() int {
	year, _ := strconv.Atoi(fmt.Sprintf("%s%s", date.Ocent(), date.Oyear()))
	month, _ := strconv.Atoi(date.Omonth())
	day, _ := strconv.Atoi(date.Oday())

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	_, w := t.ISOWeek()
	return w
}

//Wday - returns a day of week from odate
func (date Odate) Wday() int {
	year, _ := strconv.Atoi(fmt.Sprintf("20%s", date.Oyear()))
	month, _ := strconv.Atoi(date.Omonth())
	day, _ := strconv.Atoi(date.Oday())

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return int(t.Weekday())
}

//CurrentOdate - returns current odate in local time
func CurrentOdate() Odate {
	t := time.Now()
	y, m, d := t.Date()
	odat := fmt.Sprintf("%s%02s%02s", strconv.Itoa(y), strconv.Itoa(int(m)), strconv.Itoa(d))

	return Odate(odat)
}

//IsInDayOfWeek - check if tasks day of week is in odate
func IsInDayOfWeek(odate Odate, values []taskdef.ExecutionValue) bool {

	wday := odate.Wday()
	for _, val := range values {

		ival, _ := strconv.Atoi(string(val))
		if ival == wday {
			return true
		}
	}

	return false
}

//IsInDayOfMonth - check if day of execution is in odate
func IsInDayOfMonth(odate Odate, values []taskdef.ExecutionValue) bool {

	day := odate.Oday()
	for _, val := range values {
		if string(val) == day {
			return true
		}
	}

	return false
}

//IsInExactDate check if tasks execiution date is in odate
func IsInExactDate(odate Odate, values []taskdef.ExecutionValue) bool {

	dt := odate.FormatDate()
	for _, val := range values {

		if string(val) == dt {
			return true
		}
	}

	return false
}

//IsInMonth - check if tasks month is in odate
func IsInMonth(odate Odate, values []taskdef.MonthData) bool {

	mth, _ := strconv.Atoi(odate.Omonth())
	for _, val := range values {

		if int(val) == mth {
			return true
		}

	}

	return false
}

//IsBeforeCurrent - checks if order date is before current day
func IsBeforeCurrent(odate Odate, currentOdate Odate) bool {

	y, m, d := odate.Ymd()
	y2, m2, d2 := currentOdate.Ymd()
	otime := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	otime2 := time.Date(y2, time.Month(m2), d2, 0, 0, 0, 0, time.Local)

	return otime.Before(otime2)
}

//AddDays - adds num days to given date and return new odate
func AddDays(odate Odate, num int) (Odate, error) {

	y, m, d := odate.Ymd()
	otime := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

	ny, nm, nd := otime.AddDate(0, 0, num).Date()
	odat := fmt.Sprintf("%s%02s%02s", strconv.Itoa(ny), strconv.Itoa(int(nm)), strconv.Itoa(nd))

	return Odate(odat), nil

}

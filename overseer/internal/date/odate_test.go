package date

import (
	"fmt"
	"goscheduler/overseer/internal/taskdef"

	"strconv"
	"testing"
	"time"
)

func TestOdate(t *testing.T) {

	ctime := time.Now()
	result := CurrentOdate()
	year := strconv.Itoa(ctime.Year())
	month := fmt.Sprintf("%02d", ctime.Month())
	day := fmt.Sprintf("%02d", ctime.Day())
	dyear := ctime.YearDay()
	wday := int(ctime.Weekday())
	_, wofyear := ctime.ISOWeek()

	tOdat := fmt.Sprintf("%s%s%s", year, month, day)

	y, m, d := result.Ymd()
	if y != ctime.Year() && m != int(ctime.Month()) && d != ctime.Day() {
		t.Error("Ymd failed")
	}

	if tOdat != string(result) {
		t.Error("invalid odate expected:", tOdat, " actual:", result)
	}

	if result.ODATE() != tOdat[2:] {
		t.Error("invalid Odate expected:", tOdat[2:], " actual:", result.ODATE())
	}

	if result.Oday() != day {
		t.Error("invalid Oday expected:", day, " actual:", result.Oday())
	}
	if result.Omonth() != month {
		t.Error("invalid Omonth expected:", month, " actual:", result.Omonth())
	}

	if result.Ocent() != year[0:2] {
		t.Error("invalid Oyear expected:", year[0:2], " actual:", result.Ocent())
	}

	if result.Oyear() != year[2:] {
		t.Error("invalid Oyear expected:", year[2:], " actual:", result.Oyear())
	}

	if result.Doyear() != dyear {
		t.Error("invalid Doyear expected:", dyear, " actual:", result.Doyear())
	}
	if result.Wday() != wday {
		t.Error("invalid Wday expected:", wday, " actual:", result.Wday())
	}
	if result.Woyear() != wofyear {
		t.Error("invalid Woyear expected:", wofyear, " actual:", result.Woyear())
	}

	form := fmt.Sprintf("%s-%s-%s", year, month, day)
	if result.FormatDate() != form {
		t.Error("invalid Format expected:", form, " actual:", result.FormatDate())
	}

}
func TestDateRange(t *testing.T) {

	odate := Odate("20200905")
	var res bool = false
	res = IsInDayOfMonth(odate, []taskdef.ExecutionValue{"05"})
	if res == false {
		t.Error("odate not in day of month")
	}
	res = IsInDayOfMonth(odate, []taskdef.ExecutionValue{"06"})
	if res {
		t.Error("odate in day of month")
	}

	res = IsInExactDate(odate, []taskdef.ExecutionValue{"2020-09-05"})
	if res == false {
		t.Error("odate not equal exact date")
	}
	res = IsInExactDate(odate, []taskdef.ExecutionValue{"2020-09-06"})
	if res == true {
		t.Error("odate equal exact date")
	}

	res = IsBeforeCurrent(odate, Odate("20200906"))
	if res == false {
		t.Error("odate is not before current")
	}

	res = IsBeforeCurrent(odate, Odate("20200904"))
	if res == true {
		t.Error("odate is before current")
	}

	res = IsInMonth(odate, []taskdef.MonthData{8, 9})
	if res == false {
		t.Error("odate is not in month")
	}

	res = IsInMonth(odate, []taskdef.MonthData{7, 10})
	if res == true {
		t.Error("odate is in month")
	}

	res = IsInDayOfWeek(odate, []taskdef.ExecutionValue{"1", "6"})
	if res == false {
		t.Error("odate is not in day of week")
	}
	res = IsInDayOfWeek(odate, []taskdef.ExecutionValue{"0", "5"})
	if res == true {
		t.Error("odate is in day of week")
	}

	ndat, _ := AddDays(odate, 10)
	if string(ndat) != "20200915" {
		t.Error(fmt.Sprintf("Add day unexpected value:%s,expeted:%s", ndat, "20200915"))
	}

}

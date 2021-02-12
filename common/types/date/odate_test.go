package date

import (
	"fmt"

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
	var res bool
	res = IsInDayOfMonth(odate, []int{5})
	if res == false {
		t.Error("odate not in day of month")
	}
	res = IsInDayOfMonth(odate, []int{6})
	if res {
		t.Error("odate in day of month")
	}

	res = IsInExactDate(odate, []string{"2020-09-05"})
	if res == false {
		t.Error("odate not equal exact date")
	}
	res = IsInExactDate(odate, []string{"2020-09-06"})
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

	res = IsInMonth(odate, []time.Month{8, 9})
	if res == false {
		t.Error("odate is not in month")
	}

	res = IsInMonth(odate, []time.Month{7, 10})
	if res == true {
		t.Error("odate is in month")
	}

	res = IsInDayOfWeek(odate, []int{1, 6})
	if res == false {
		t.Error("odate is not in day of week")
	}
	res = IsInDayOfWeek(odate, []int{0, 5})
	if res == true {
		t.Error("odate is in day of week")
	}

	ndat := AddDays(odate, 10)
	if string(ndat) != "20200915" {
		t.Error(fmt.Sprintf("Add day unexpected value:%s,expeted:%s", ndat, "20200915"))
	}

}

func TestValidate(t *testing.T) {

	ok, err := Odate("").validateValue()

	if !ok {
		t.Error(err)
	}

	//Tests some valid values
	tvalues := []string{"00010101", "99991231", "20000229", "29990101", "19840715", "20200330", "20200331", "20200531"}

	for _, value := range tvalues {
		ok, err := Odate(value).validateValue()

		if !ok {
			t.Error(err, value)
		}
	}

	//Tests length of an odate
	ok, err = Odate("0123456").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidLen {
		t.Error("unexpected value, epxected:", errOdateInvalidLen, "actual", err)
	}

	ok, err = Odate("012345678").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidLen {
		t.Error("unexpected value, epxected:", errOdateInvalidLen, "actual", err)
	}

	//Tests if an odate contains only numeric values

	ok, err = Odate("2222111a").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateNotNumeric {
		t.Error("unexpected value, epxected:", errOdateNotNumeric, "actual", err)
	}

	//Tests a range of year value
	ok, err = Odate("00001115").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidYear {
		t.Error("unexpected value, epxected:", errOdateInvalidYear, "actual", err)
	}

	ok, err = Odate("20200000").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidMonth {
		t.Error("unexpected value, epxected:", errOdateInvalidMonth, "actual", err)
	}

	//Tests a range of month value
	ok, err = Odate("20201300").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidMonth {
		t.Error("unexpected value, epxected:", errOdateInvalidMonth, "actual", err)
	}

	ok, err = Odate("20200100").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidDay {
		t.Error("unexpected value, epxected:", errOdateInvalidDay, "actual", err)
	}

	//Tests range of day value
	ok, err = Odate("20200132").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidDay {
		t.Error("unexpected value, epxected:", errOdateInvalidDay, "actual", err)
	}

	//Tests end of month value
	fmonths := []string{"20200230", "20200431", "20200631", "20200931", "20201131"}

	for _, value := range fmonths {
		ok, err := Odate(value).validateValue()

		if ok == true {
			t.Error("unexpected value, epxected:", false)
		}
		if err != errOdateInvalidDay {
			t.Error("unexpected value, epxected:", errOdateInvalidDay, "actual", err)
		}
	}

	//Tests leap year
	tlyears := []string{"20200229", "20000229"}

	for _, value := range tlyears {
		ok, err := Odate(value).validateValue()

		if ok != true {
			t.Error(err)
		}
	}

	ok, err = Odate("19000229").validateValue()
	if ok == true {
		t.Error("unexpected value, epxected:", false)
	}
	if err != errOdateInvalidDay {
		t.Error("unexpected value, epxected:", errOdateInvalidDay, "actual", err)
	}

}
func TestFromTime(t *testing.T) {

	r := FromTime(time.Now())
	if r != CurrentOdate() {
		t.Error("unexpected result:")
	}
}

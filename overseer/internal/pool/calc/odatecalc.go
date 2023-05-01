package calc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/pool/models"
	"github.com/przebro/overseer/overseer/internal/taskdef"
)

type ExecutedTaskInfo interface {
	GetInfo() (string, string, string)
	ExecutionID() string
	RunNumber() int32
	OrderID() unique.TaskOrderID
	Variables() types.EnvironmentVariableList
}

func CalcRealOdate(current date.Odate, expect date.OdateValue, schedule taskdef.SchedulingData) date.Odate {

	mths := map[time.Month]bool{}
	for i := 1; i < 13; i++ {
		mths[time.Month(i)] = false
	}

	for _, n := range schedule.Months {
		mths[n] = true
	}

	//if expected odate value is any
	if expect == date.OdateValueAny || expect == date.OdateValueNone {
		return date.OdateNone
	}

	//For all ordering types, expect = ODATE means current day
	if expect == date.OdateValueDate {
		return current
	}

	//It is explicite date relative to current date, so do just simply compute
	if expect != date.OdateValueNext && expect != date.OdateValuePrev {

		days, _ := strconv.Atoi(string(expect))
		return date.AddDays(current, days)
	}

	//Manual ordering means that no specific schedule criteria being used, NEXT and PREV means simply tomorrow and yesterday
	if schedule.OrderType == taskdef.OrderingManual {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}

		return date.AddDays(current, days)
	}

	var result date.Odate

	/*for daily ordering in a simple scenario means tomorrow day for NEXT and yesterday for PREV, however,
	if PREV resolves to a day in the previous month and the execution of a task is excluded in that month then PREV will resolve
	to a day in a first included month, before the current one.
	for instance, if the task can run in months:[4,5,7]
	PREV for 05-01 means 04-30, but PREV for 07-01 means 05-31
	in particular case where months:[10], PREV for 2020-10-01 means 2019-10-31.
	The same rules apply to NEXT.
	*/
	if schedule.OrderType == taskdef.OrderingDaily {

		result = calcDateDaily(current, expect, mths)
	}

	//for exact ordering, NEXT and PREV means neighbour value in the array of specified dates
	//for edge cases(single value,last value or first value) corresponding NEXT and PREV resolves to +1/-1
	//If a task is forced to run on a non scheduled day then NEXT and PREV will resolve to +1/-1
	//for instance: if execdates: [2020-05-11] and the order date is 20200501 then NEXT will resolve to 20200511 and PREV to 20200430
	if schedule.OrderType == taskdef.OrderingExact {

		result = calcDateExact(current, expect, schedule.Exactdates)
	}

	//for weekly ordering NEXT and PREV means next or previous value from specified days of a week.
	//if it is the last specified day then the first day of a next week will be used as a value for NEXT
	//respectively, if it is the first day, then the last day from a previous week will be used for PREV
	//for instance: when a task is ordered in the day of week [1 3 4 6]
	//if it is the fourth day of a week then NEXT mean 6 and PREV means 1, however,
	//if it is the first day NEXT means 3 but PREV means 6
	if schedule.OrderType == taskdef.OrderingWeek {
		result = calcDateWeek(current, expect, schedule.Dayvalues, mths)

	}

	//day of month. Task is ordered on specific day
	// if the task is ordered on the day of the month [31] when the date is 2020-03-31 NEXT means 2020-05-31 and PREV means 2020-01-30
	//because there is no such date like 2020-02-30 and 2020-04-31

	if schedule.OrderType == taskdef.OrderingDayOfMonth {

		result = calcDateMonth(current, expect, schedule.Dayvalues, mths)
	}

	// From end of month, where 1 means last day of month,2 means a day before last day
	if schedule.OrderType == taskdef.OrderingFromEnd {

		result = calcDateFromEnd(current, expect, schedule.Dayvalues, mths)
	}

	return result
}

func calcDateDaily(current date.Odate, expect date.OdateValue, mths map[time.Month]bool) date.Odate {

	days := 1
	if expect == date.OdateValuePrev {
		days *= -1
	}

	planed := current

	for {
		planed = date.AddDays(planed, days)
		py, pm, day := planed.Ymd()
		if mths[time.Month(pm)] {
			planed = date.Odate(fmt.Sprintf("%d%02d%02d", py, pm, day))
			break
		}
	}

	return planed
}

func calcDateExact(current date.Odate, expect date.OdateValue, dates []string) date.Odate {

	var idx int
	var val string
	var found bool
	cdat := current.FormatDate()

	//first check if the task was forced in current day
	for idx, val = range dates {
		if val == cdat {
			found = true
			break
		}
	}
	// task was forced on a non scheduled day or it is only single value so return -1 or +1
	if !found || len(dates) == 1 {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	//edge case for PREV(it is first the first execution) and NEXT(it is the last execution)
	if (idx == 0 && expect == date.OdateValuePrev) || (idx == len(dates)-1 && expect == date.OdateValueNext) {

		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	//for any other case it is the next or previous item from the table of values
	nextval := 1
	if expect == date.OdateValuePrev {
		nextval *= -1
	}
	return date.FromDateString(dates[idx+nextval])

}

func calcDateFromEnd(current date.Odate, expect date.OdateValue, values []int, mths map[time.Month]bool) date.Odate {

	var shift = values[0]

	cm, cy := getNextMonthYear(mths, current, expect, false)

	d := getNthLastDay(cy, cm, shift)

	return date.Odate(fmt.Sprintf("%d%02d%02d", cy, cm, d))

}

func calcDateWeek(current date.Odate, expect date.OdateValue, values []int, mths map[time.Month]bool) date.Odate {

	var idx int
	var val int
	var diffWeek int
	var expectVal int
	var found bool

	cdat := current.Wday()

	for idx, val = range values {
		if val == cdat {
			found = true
			break
		}
	}

	//Task was forced on a no scheduled day
	if !found {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	nval := 1
	if expect == date.OdateValuePrev {
		nval *= -1
	}

	refdate := current

	for {

		idx += nval

		if idx < 0 && expect == date.OdateValuePrev {
			idx = len(values) - 1
			expectVal = values[idx]
			diffWeek = -1
		} else if idx > (len(values)-1) && expect == date.OdateValueNext {
			idx = 0
			expectVal = values[idx]
			diffWeek = 1
		} else {
			expectVal = values[idx]
			diffWeek = 0
		}

		refdate = getStartOfWeek(refdate, diffWeek)
		refdate = date.AddDays(refdate, expectVal-1)
		_, cm, _ := refdate.Ymd()

		if mths[time.Month(cm)] {
			break
		}
	}

	return refdate
}

func calcDateMonth(current date.Odate, expect date.OdateValue, values []int, mths map[time.Month]bool) date.Odate {

	var expectVal int
	var diffMonth bool
	var ndate string

	var idx int
	var val int
	var found bool

	cdat := current.Day()
	for idx, val = range values {
		if val == cdat {
			found = true
			break
		}
	}

	//Task was forced on a no scheduled day
	if !found {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	nval := 1
	if expect == date.OdateValuePrev {
		nval *= -1
	}

	refdate := current

	for {

		idx += nval

		if idx < 0 && expect == date.OdateValuePrev {
			idx = len(values) - 1
			expectVal = values[idx]
			diffMonth = true
		} else if idx > (len(values)-1) && expect == date.OdateValueNext {
			idx = 0
			expectVal = values[idx]
			diffMonth = true
		} else {
			expectVal = values[idx]
			diffMonth = false
		}

		cm, cy := getNextMonthYear(mths, refdate, expect, !diffMonth)
		lday := getNthLastDay(cy, cm, 1)

		if expectVal <= lday {
			ndate = fmt.Sprintf("%d%02d%02d", cy, cm, expectVal)
			break
		}

		refdate = date.Odate(fmt.Sprintf("%d%02d%02d", cy, cm, lday))
	}

	return date.Odate(ndate)

}

func getNextMonthYear(mths map[time.Month]bool, current date.Odate, expect date.OdateValue, incl bool) (int, int) {

	cy, cm, _ := current.Ymd()
	nval := 1
	if expect == date.OdateValuePrev {
		nval *= -1
	}

	if !incl {
		cm += nval
	}

	if cm < 1 {
		cm = 12
		cy--
	}
	if cm > 12 {
		cm = 1
		cy++
	}

	for !mths[time.Month(cm)] {
		cm += nval

		if cm < 1 {
			cm = 12
			cy--
		}
		if cm > 12 {
			cm = 1
			cy++
		}
	}

	return cm, cy
}

// getNthLastDay - gets the nth last day from given year and month
func getNthLastDay(year int, month int, shift int) int {

	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	pd := t.AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

	if shift > 1 {
		pd -= (shift - 1)
	}

	return pd
}

// getStartOfWeek - gets an odate of first day(monday) in  week
func getStartOfWeek(current date.Odate, shift int) date.Odate {
	y, m, d := current.Ymd()
	wday := current.Wday()

	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local).AddDate(0, 0, -wday+1).AddDate(0, 0, shift*7)
	return date.FromTime(t)
}

func PrepareVaribles(task ExecutedTaskInfo, odate date.Odate) types.EnvironmentVariableList {

	variables := types.EnvironmentVariableList{}
	n, _, _ := task.GetInfo()
	cID := task.ExecutionID()

	variables = append(variables, types.EnvironmentVariable{Name: "%%ORDERID", Value: string(task.OrderID())})
	variables = append(variables, types.EnvironmentVariable{Name: "%%RN", Value: fmt.Sprintf("%d", task.RunNumber())})
	variables = append(variables, types.EnvironmentVariable{Name: "%%EXECID", Value: cID})
	variables = append(variables, types.EnvironmentVariable{Name: "%%ODATE", Value: odate.ODATE()})
	variables = append(variables, types.EnvironmentVariable{Name: "%%TASKNAME", Value: n})

	variables = append(variables, task.Variables()...)

	return variables
}

func ComputeTaskState(taskType types.TaskType, maxrc, rc, sc int32) models.TaskState {

	var state models.TaskState = models.TaskStateEndedOk

	switch taskType {
	case types.TypeDummy:
		{
			state = models.TaskStateEndedOk
		}
	case types.TypeOs:
		{
			if rc > maxrc || rc < 0 {
				state = models.TaskStateEndedNotOk
			}
		}
	case types.TypeAws:
		{
			if sc == int32(types.StatusCodeNormal) {
				state = models.TaskStateEndedOk
			} else {
				state = models.TaskStateEndedNotOk
			}
		}
	}
	return state
}

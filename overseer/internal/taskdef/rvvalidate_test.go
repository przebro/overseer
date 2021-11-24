package taskdef

import (
	"overseer/common/types"
	"overseer/common/validator"
	"testing"
	"time"
)

func TestInTicketData(t *testing.T) {

	data := InTicketData{Name: "ABCDEF", Odate: ""}
	err := validator.Valid.Validate(data)
	if err != nil {
		t.Error(err)
	}

	data = InTicketData{Name: "ABCDEF", Odate: "*"}
	err = validator.Valid.Validate(data)
	if err != nil {
		t.Error(err)
	}

	data = InTicketData{Name: "ABCDEF", Odate: "ODATE"}
	err = validator.Valid.Validate(data)
	if err != nil {
		t.Error(err)
	}
	data = InTicketData{Name: "ABCDEF", Odate: "NEXT"}
	err = validator.Valid.Validate(data)
	if err != nil {
		t.Error(err)
	}
	data = InTicketData{Name: "ABCDEF", Odate: "PREV"}
	err = validator.Valid.Validate(data)
	if err != nil {
		t.Error(err)
	}

	data = InTicketData{Name: "ABCDEF", Odate: "2020"}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}

	data = InTicketData{Name: "ABCDEF", Odate: "_"}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}

	data = InTicketData{Name: "a-BcdE_123", Odate: ""}
	err = validator.Valid.Validate(data)
	if err != nil {
		t.Error(err)
	}

	data = InTicketData{Name: "!ABC", Odate: ""}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}

	data = InTicketData{Name: "_ABC", Odate: ""}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}

	data = InTicketData{Name: "0ABC", Odate: ""}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}

	data = InTicketData{Name: "!ABC", Odate: ""}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}

}

func TestOutTicketData(t *testing.T) {

	data := OutTicketData{Name: "ABC", Odate: "", Action: "ADD"}
	err := validator.Valid.Validate(data)
	if err != nil {
		t.Error("Unexpected error")
	}

	data = OutTicketData{Name: "ABC", Odate: "", Action: "REM"}
	err = validator.Valid.Validate(data)
	if err != nil {
		t.Error("Unexpected error")
	}

	data = OutTicketData{Name: "ABC", Odate: "", Action: ""}
	err = validator.Valid.Validate(data)
	if err == nil {
		t.Error("Unexpected error")
	}
}

func TestFlagData(t *testing.T) {

	flag := FlagData{Name: "ABCDEF", Type: "SHR"}
	err := validator.Valid.Validate(flag)
	if err != nil {
		t.Error(err)
	}

	flag = FlagData{Name: "ABCDEF", Type: "EXL"}
	err = validator.Valid.Validate(flag)
	if err != nil {
		t.Error(err)
	}

	flag = FlagData{Name: "ABCDEF", Type: ""}
	err = validator.Valid.Validate(flag)
	if err == nil {
		t.Error(err)
	}
}

func TestVariableData(t *testing.T) {

	variable := VariableData{Name: "%%ABCD", Value: ""}

	err := validator.Valid.Validate(variable)
	if err != nil {
		t.Error(err)
	}

	variable = VariableData{Name: "%ABCD", Value: ""}

	err = validator.Valid.Validate(variable)
	if err == nil {
		t.Error(err)
	}

	variable = VariableData{Name: "%aAAA", Value: ""}

	err = validator.Valid.Validate(variable)
	if err == nil {
		t.Error(err)
	}
}

func TestSchedulingData(t *testing.T) {

	tm := SchedulingData{
		OrderType: "manual",
		FromTime:  "23:50",
	}
	err := validator.Valid.Validate(tm)
	if err != nil {
		t.Error(err)
	}

	tm = SchedulingData{
		OrderType: "manual",
		FromTime:  "24:50",
	}
	err = validator.Valid.Validate(tm)
	if err == nil {
		t.Error("Unexpected value")
	}

	tm = SchedulingData{
		OrderType: "manual",
		FromTime:  "",
	}
	err = validator.Valid.Validate(tm)
	if err != nil {
		t.Error(err)
	}

	tm = SchedulingData{
		OrderType: "manual",
		Months:    []time.Month{},
	}
	err = validator.Valid.Validate(tm)
	if err != nil {
		t.Error(err)
	}

	tm = SchedulingData{
		OrderType: "manual",
		Months:    []time.Month{time.January, time.January},
	}
	err = validator.Valid.Validate(tm)
	if err == nil {
		t.Error("Unexpected value")
	}

}

func TestTaskDefinition(t *testing.T) {

	sd := SchedulingData{
		OrderType: "manual",
	}

	tdf := baseTaskDefinition{
		TaskType:  types.TaskType("dummy"),
		Name:      "AAA",
		Group:     "ABC",
		Schedule:  sd,
		InTickets: []InTicketData{{Name: "_____"}},
	}
	err := validator.Valid.Validate(tdf)

	if err == nil {
		t.Error("unexpected value")
	}

	tdf = baseTaskDefinition{
		TaskType:   types.TaskType("dummy"),
		Name:       "AAA",
		Group:      "ABC",
		Schedule:   sd,
		InTickets:  []InTicketData{},
		InRelation: "AND",
	}
	err = validator.Valid.Validate(tdf)

	if err != nil {
		t.Error(err)
	}
}

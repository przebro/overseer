package taskdef

import (
	"fmt"
	"io/ioutil"
	"os"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/common/types/date"
	"overseer/common/validator"
	"overseer/overseer/taskdata"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var expect TaskDefinition = &baseTaskDefinition{
	Name: "dummy_01", Group: "", Description: "sample dummy task definition", ConfirmFlag: false, TaskType: "dummy",
	InTickets:  []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
	InRelation: InTicketAND,
	OutTickets: []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
	FlagsTab:   []FlagData{{Name: "flag01", Type: FlagShared}},
	Schedule: SchedulingData{
		OrderType: "weekday",
		FromTime:  "11:30",
		ToTime:    "",
		Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Dayvalues: []int{1, 3, 5},
	},
}

var expect3 TaskDefinition = &baseTaskDefinition{
	Name: "dummy_04", Group: "test", Description: "sample modified dummy task definition", ConfirmFlag: false, TaskType: types.TypeOs,
	DataRetention: 1,
	InTickets:     []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
	InRelation:    InTicketAND,
	OutTickets:    []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
	FlagsTab:      []FlagData{{Name: "flag01", Type: FlagExclusive}},
	Schedule: SchedulingData{
		OrderType:    OrderingDayOfMonth,
		FromTime:     "15:30",
		ToTime:       "",
		Months:       []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Dayvalues:    []int{1, 3, 5},
		AllowPastSub: false,
	},
}

var expectOutTicket string = `"outticket":[{"name":"OK-COND-02","odate":"ODATE","action":"ADD"}]`

func TestInit(t *testing.T) {

}

func TestLoad(t *testing.T) {
	t.Log(expect)
}

func TestTaskData(t *testing.T) {
	if expect3.OrderType() != OrderingDayOfMonth {
		t.Error("task data invalid order type")
	}
	if expect3.AllowPast() != false {
		t.Error("task data invalid allow past ")
	}
	if expect3.Confirm() != false {
		t.Error("task data invalid  confirmflag")
	}
	if expect3.Retention() != 1 {
		t.Error("task data invalid  retention")
	}
	if len(expect3.Months()) != 12 {
		t.Error("task data invalid  months")
	}
	if expect3.TypeName() != types.TypeOs {
		t.Error("task data invalid  tasktype")
	}
}

func TestUnmarshalTask(t *testing.T) {

	pth, _ := filepath.Abs(filepath.Join(defDirectory, `test/dummy_01.json`))

	result, err := FromDefinitionFile(pth)
	if err != nil {
		t.Log(pth)
		t.Error("unable  to deserialize definition")
	}

	from, _ := expect.TimeSpan()
	rfrom, _ := result.TimeSpan()
	if from != rfrom {
		t.Errorf("Unmarshal failed, FromTime not equal")

	}
	if len(expect.Days()) != len(result.Days()) {
		t.Error("Unmarshal failed, Values not equal")
	}
	for x, n := range expect.Days() {
		if result.Days()[x] != n {
			t.Error("Unmarshal failed, Values not equal")
		}
	}
}
func TestMarshalTask(t *testing.T) {

	data, err := SerializeDefinition(expect)
	if err != nil {
		t.Error("Marshal failed:", err)
	}

	dstr := data

	pos := strings.Index(dstr, `"outticket":`)
	if pos == -1 {
		t.Fatal("Marshal, unable to find given substring")
	}
	substr := dstr[pos : pos+len(expectOutTicket)]

	if substr != expectOutTicket {
		t.Error("Marshal, compared substrings does not match")
	}

	pos = strings.Index(dstr, `"schedule":`)
	if pos == -1 {
		t.Fatal("Marshal, unable to find given substring")
	}

}

func TestManager_Create_Task(t *testing.T) {

	path, _ := filepath.Abs(defDirectory)
	manager, _ := NewManager(path, logger.NewTestLogger())
	_, modTask2, _, _ := helperCreateTasks()

	//Create new task
	err := manager.Create(modTask2)
	if err != nil {
		t.Error("Create error", err)

	}
}

func TestManager_Create_Error(t *testing.T) {

	path, _ := filepath.Abs(defDirectory)
	manager, _ := NewManager(path, logger.NewTestLogger())
	modTask, _, _, _ := helperCreateTasks()
	//Try create task that already exists
	err := manager.Create(modTask)
	if !strings.Contains(err.Error(), "unable to create, definition already exists") {
		t.Error("Create error", err)

	}
}

func TestManager_Delete(t *testing.T) {

	path, _ := filepath.Abs(defDirectory)
	manager, _ := NewManager(path, logger.NewTestLogger())
	_, modTask2, _, _ := helperCreateTasks()

	name, grp, _ := modTask2.GetInfo()

	//Create new task ignore error
	manager.Create(modTask2)

	err := manager.Delete(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: grp}, Name: name})
	if err != nil {
		t.Error("Create error", err)

	}
}

func TestManager_Delete_Error(t *testing.T) {

	path, _ := filepath.Abs(defDirectory)
	manager, _ := NewManager(path, logger.NewTestLogger())
	err := manager.Delete(taskdata.GroupNameData{Name: "dummy_AA", GroupData: taskdata.GroupData{Group: "test"}})
	if err == nil {
		t.Error("Delete error")
	}

	if _, ok := err.(*os.PathError); !ok {
		t.Error("Delete error")
	}
}

func TestGetTask(t *testing.T) {

	manager, err := NewManager(managerPath, logger.NewTestLogger())
	if err != nil {
		t.Fatal("unable to intialize manager")
	}

	tlist := []taskdata.GroupNameData{{Name: "dummy_01", GroupData: taskdata.GroupData{Group: "test"}},
		{Name: "task_that_does_not_exists", GroupData: taskdata.GroupData{Group: "test"}}}

	result := manager.GetTasks(tlist...)
	if len(result) != 1 {
		t.Error("unexpected result, expected :", len(tlist), "got:", len(result))
	}

	_, err = manager.GetTasksFromGroup([]string{"test", "no_group_name"})
	if err == nil {
		t.Error("unexpected value,group does not exists")
	}

	_, err = manager.GetTasksFromGroup([]string{"test"})
	if err != nil {
		t.Error("unexpected value,group exists")
	}

}

func TestListTaskModel(t *testing.T) {
	path, _ := filepath.Abs(defDirectory)
	manager, err := NewManager(path, logger.NewTestLogger())
	if err != nil {
		t.Fatal("unable to intialize manager")
	}

	if _, err := manager.GetTaskModelList(taskdata.GroupData{Group: "dir_not_exists"}); err == nil {
		t.Error("unexpected result")
	}

	list, err := manager.GetTaskModelList(taskdata.GroupData{Group: "test"})
	if err != nil {
		t.Error("unexpected result:", err)
	}

	dirpath := filepath.Join(path, "test")
	info, err := ioutil.ReadDir(dirpath)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if len(info) != len(list) {
		t.Error("unexpected result: directory list is not same as task model list")
	}

	for i, n := range info {

		if strings.HasPrefix(n.Name(), list[i].Name) != true {
			t.Error("unexpected result, task name not equal file name:", list[i].Name, n.Name())
		}
	}
}

func TestMarshalTests2(t *testing.T) {

	data, _ := SerializeDefinition(expect3)
	_, err := FromString(data)
	if err != nil {
		t.Error("Unmarshal error")
	}

	_, err = FromString(`{"type" : "dummy","name" :"sample_01A","group" : "samples"`)
	if !strings.Contains(err.Error(), "unexpected end of JSON") {
		t.Error(err)
	}

	_, err = FromString(`{"type" : "dummy","name" :"","group" : "samples"}`)
	if err == nil {
		t.Error("Unexpected result")
	}

}

func TestGetTimeSpan(t *testing.T) {

	schdata := SchedulingData{FromTime: "10:30", ToTime: "11:30", OrderType: OrderingManual}

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_time_span", "description").
		WithSchedule(schdata).
		WithFlags([]FlagData{{Name: "FLAG01", Type: FlagShared}}).
		WithConfirm().WithRetention(1).WithVariables([]VariableData{{Name: "%%VAR", Value: "xx"}}).Build()

	if err != nil {
		t.Error("task builder error", err)
	}

	from, to := def.TimeSpan()
	if from.String() != "10:30" || to.String() != "11:30" {
		t.Error("Unexpected values:", from, to)
	}
}

func TestCyclicData(t *testing.T) {

	path, _ := filepath.Abs(defDirectory)
	manager, err := NewManager(path, logger.NewTestLogger())
	if err != nil {
		t.Fatal("unable to intialize manager")
	}

	def, err := manager.GetTask(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "cyclic_01"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(def.GetInfo())

}

func TestValidateCyclicData(t *testing.T) {

	var ctd CyclicTaskData
	var err error
	ctd = CyclicTaskData{}

	if err = validator.Valid.Validate(ctd); err != nil {
		t.Error("unexpected result:", err)
	}

	ctd.TimeInterval = -1

	if err = validator.Valid.Validate(ctd); err == nil {
		t.Error("unexpected result:")
	}

	ctd.TimeInterval = 1441

	if err = validator.Valid.Validate(ctd); err == nil {
		t.Error("unexpected result:")
	}

	ctd.TimeInterval = 5
	ctd.MaxRuns = -1

	if err = validator.Valid.Validate(ctd); err == nil {
		t.Error("unexpected result:")
	}

	ctd.MaxRuns = 1000

	if err = validator.Valid.Validate(ctd); err == nil {
		t.Error("unexpected result:")
	}

	ctd.MaxRuns = 5
	ctd.RunFrom = "start"

	if err = validator.Valid.Validate(ctd); err != nil {
		t.Error("unexpected result:")
	}

	ctd.RunFrom = "end"

	if err = validator.Valid.Validate(ctd); err != nil {
		t.Error("unexpected result:")
	}

	ctd.RunFrom = "schedule"

	if err = validator.Valid.Validate(ctd); err != nil {
		t.Error("unexpected result:")
	}

	ctd.RunFrom = "unknown"

	if err = validator.Valid.Validate(ctd); err == nil {
		t.Error("unexpected result:")
	}

	fmt.Println(ctd)

}

func TestGetAction(t *testing.T) {

	schdata := SchedulingData{FromTime: "10:30", ToTime: "11:30", OrderType: OrderingDaily}

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_time_span", "description").
		WithSchedule(schdata).
		WithFlags([]FlagData{{Name: "FLAG01", Type: FlagShared}}).
		WithConfirm().WithRetention(1).WithVariables([]VariableData{{Name: "%%VAR", Value: "xx"}}).Build()

	if err != nil {
		t.Error("task builder error")
	}

	if string(def.Action()) != "" {
		t.Error("Unexpected value")
	}
}

func TestExpandVariable(t *testing.T) {

	expect := "OVS_VARIABLE"
	variable := VariableData{Name: "%%VARIABLE", Value: ""}
	if variable.Expand() != expect {
		t.Error("Unexpected value expected:", expect, "actual", variable.Expand())
	}

}

func TestBuilder(t *testing.T) {

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_04", "description").
		WithFlags([]FlagData{{Name: "FLAG01", Type: FlagShared}}).
		WithSchedule(SchedulingData{OrderType: OrderingDaily}).
		WithConfirm().WithRetention(1).WithVariables([]VariableData{{Name: "%%VAR", Value: "xx"}}).Build()

	if err != nil {
		t.Error("task builder error")
	}

	if def.Confirm() != true {
		t.Error("task builder error expected:", true, " actual:", def.Confirm())
	}

	if def.Retention() != 1 {
		t.Error("task builder error expected:", 1, " actual:", def.Retention())
	}

	if len(def.Variables()) != 1 {
		t.Error("task builder error expected:", 1, " actual:", len(def.Variables()))
	}

}

func helperCreateTasks() (t1, t2, t3, t4 TaskDefinition) {

	t1 = &baseTaskDefinition{
		TaskType: "dummy",
		Name:     "dummy_02", Group: "test", Description: "sample modified dummy task definition", ConfirmFlag: false,
		InTickets:  []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
		InRelation: InTicketAND,
		OutTickets: []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
		FlagsTab:   []FlagData{{Name: "flag01", Type: FlagShared}},
		Schedule: SchedulingData{
			OrderType: "weekday",
			FromTime:  "15:30",
			ToTime:    "",
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{1, 3, 5},
		},
	}

	t2 = &baseTaskDefinition{
		TaskType: "dummy",
		Name:     "dummy_AA", Group: "test", Description: "sample modified dummy task definition", ConfirmFlag: false,
		InTickets:  []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
		InRelation: InTicketAND,
		OutTickets: []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
		FlagsTab:   []FlagData{{Name: "flag01", Type: FlagShared}},
		Schedule: SchedulingData{
			OrderType: "weekday",
			FromTime:  "11:30",
			ToTime:    "",
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{1, 3, 5},
		},
	}

	t3 = &baseTaskDefinition{
		TaskType: "dummy",
		Name:     "", Group: "test", Description: "sample modified dummy task definition", ConfirmFlag: false,
		InTickets:  []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
		InRelation: InTicketAND,
		OutTickets: []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
		FlagsTab:   []FlagData{{Name: "flag01", Type: FlagShared}},
		Schedule: SchedulingData{
			OrderType: "weekday",
			FromTime:  "11:30",
			ToTime:    "",
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{1, 3, 5},
		},
	}
	t4 = &baseTaskDefinition{
		TaskType: "dummy",
		Name:     "dummy_01", Group: "test", Description: "sample modified dummy task definition", ConfirmFlag: false,
		InTickets:  []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
		InRelation: InTicketAND,
		OutTickets: []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
		FlagsTab:   []FlagData{{Name: "flag01", Type: FlagShared}},
		Schedule: SchedulingData{
			OrderType: "weekday",
			FromTime:  "11:30",
			ToTime:    "",
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{1, 3, 5},
		},
	}
	return
}

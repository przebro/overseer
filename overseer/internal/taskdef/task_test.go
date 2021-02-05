package taskdef

import (
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/overseer/internal/date"
	"overseer/overseer/taskdata"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var log logger.AppLogger = logger.NewTestLogger()

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
		Values:    []string{"1", "3", "5"},
	},
}
var expect2 TaskDefinition = &baseTaskDefinition{
	Name: "dummy_02", Group: "", Description: "sample modified dummy task definition", ConfirmFlag: false, TaskType: "dummy",
	InTickets:  []InTicketData{{Name: "OK-COND-01", Odate: date.OdateValueDate}},
	InRelation: InTicketAND,
	OutTickets: []OutTicketData{{Name: "OK-COND-02", Action: OutActionAdd, Odate: date.OdateValueDate}},
	FlagsTab:   []FlagData{{Name: "flag01", Type: FlagShared}},
	Schedule: SchedulingData{
		OrderType: "weekday",
		FromTime:  "15:30",
		ToTime:    "",
		Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Values:    []string{"1", "3", "5"},
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
		Values:       []string{"1", "3", "5"},
		AllowPastSub: false,
	},
}

var expectOutTicket string = `"outticket":[{"name":"OK-COND-02","odate":"ODATE","action":"ADD"}]`
var expectSchedule string = `"schedule":{"type":"weekday","from":"15:30","to":"","months":[1,2,3,4,5,6,7,8,9,10,11,12],"values":["1","3","5"]}`

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

	pth, _ := filepath.Abs(`../../../def/test/dummy_01.json`)

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
	if len(expect.Values()) != len(result.Values()) {
		t.Error("Unmarshal failed, Values not equal")
	}
	for x, n := range expect.Values() {
		if result.Values()[x] != n {
			t.Error("Unmarshal failed, Values not equal")
		}
	}
}
func TestMarshalTask(t *testing.T) {

	data, err := WriteDefinitionFile(expect)
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

func TestManagerGroups(t *testing.T) {
	path, _ := filepath.Abs("../../../def/")
	manager, _ := NewManager(path)

	groups := manager.GetGroups()

	if groups[0] != "" && groups[1] != "test" {
		t.Error("invalig group names")
	}

	err := manager.DeleteGroup("")
	if err.Error() != "group name cannot be empty" {
		t.Error("Unexpected error")
	}

	err = manager.DeleteGroup("test")
	if err == nil || err.Error() != "directory is not empty" {
		t.Error("Manager, non empty group should not be deleted.", err)
	}

	err = manager.DeleteGroup("test2")
	if err == nil || !strings.Contains(err.Error(), "can't find directory") {
		t.Error("Delete group", err)
	}

	err = manager.CreateGroup("test2")
	if err != nil {
		t.Error("Error, create new group")
	}

	err = manager.CreateGroup("test2")
	if err == nil {
		t.Error("Error, create new group,already exists")
	}

	err = manager.CreateGroup("")
	if err == nil {
		t.Error("Error, create new group. group name is empty")
	}

	manager.DeleteGroup("test2")

}

func TestManagerLock(t *testing.T) {

	path, _ := filepath.Abs("../../../def/")
	manager, _ := NewManager(path)
	// try to unlock a lock that does not exists
	err := manager.Unlock(12345)
	if err == nil {
		t.Log("Unlock error, nonexistent value")
	}

	//locking an empty task is not possible
	_, err = manager.Lock(taskdata.GroupNameData{Group: "", Name: ""})
	if err != ErrTaskNameEmpty {
		t.Error("Lock error,empty task name not allowed")
	}

	// acquire new lock
	lID, err := manager.Lock(taskdata.GroupNameData{Group: "test", Name: "dummy_01"})
	if err != nil {
		t.Error("Lock error")
	}

	//acquiring a new lock for a task that is already locked isn't possible
	_, err = manager.Lock(taskdata.GroupNameData{Group: "test", Name: "dummy_01"})
	if !strings.Contains(err.Error(), "Unable to acquire lock") {
		t.Error("Lock error, object already locked")
	}

	// locking a task that doesn't exisits isn't possible
	_, err = manager.Lock(taskdata.GroupNameData{Group: "test", Name: "dummy_0123"})
	if err == nil {
		t.Error("Lock error,a file that not exists locked")
	}

	// release the lock
	err = manager.Unlock(lID)
	if err != nil {
		t.Error("Unlock error")
	}

}
func TestManagerUpdate(t *testing.T) {

	path, _ := filepath.Abs("../../../def/")

	manager, _ := NewManager(path)
	modTask, modTask2, modTask3, modTask4 := helperCreateTasks()

	lockID, err := manager.Lock(taskdata.GroupNameData{Name: "dummy_01", Group: "test"})
	if err != nil {
		t.Fatal("Unable to lock task", err)
	}
	def := manager.GetTasks(taskdata.GroupNameData{Name: "dummy_01", Group: "test"})

	if len(def) == 0 {
		t.Fatal("def expected not empty")
	}

	if err != nil {
		t.Fatal("Unable to acquire lock")
	}

	//Try update task without associated lockId
	err = manager.Update(0, modTask)

	if !strings.Contains(err.Error(), "given lockID does not exists") {
		t.Log("Update error,invalid taskID")
	}
	//Try override definition
	err = manager.Update(lockID, modTask)

	if err != nil && !strings.Contains(err.Error(), "unable to rename,") {
		t.Error("Update error,task ovveride", err)
	}
	//Only rename to new name is possible
	err = manager.Update(lockID, modTask2)
	if err != nil {
		t.Error("Update task, task rename error", err)
	}

	//Try create task without name
	err = manager.Update(lockID, modTask3)
	if !strings.Contains(err.Error(), "task name cannot be empty") {
		t.Error("Update task, empty name", err)
	}

	//Back to previous name
	err = manager.Update(lockID, modTask4)
	if err != nil {
		t.Error("Update task, task rename error", err)
	}
}

func TestManagerCreateDelete(t *testing.T) {

	path, _ := filepath.Abs("../../../def/")
	manager, _ := NewManager(path)
	modTask, modTask2, _, _ := helperCreateTasks()

	name, grp, _ := modTask.GetInfo()
	name2, grp2, _ := modTask2.GetInfo()

	err := manager.Delete(0, taskdata.GroupNameData{Group: grp, Name: name})
	if !strings.Contains(err.Error(), "given lockID does not exists") {
		t.Error("Delete error", err)

	}

	lockID, err := manager.Lock(taskdata.GroupNameData{Name: "dummy_01", Group: "test"})

	err = manager.Delete(lockID, taskdata.GroupNameData{Name: "", Group: ""})
	if !strings.Contains(err.Error(), "group and name does not match with lockID") {
		t.Error("Delete error", err)

	}
	err = manager.Delete(lockID, taskdata.GroupNameData{Name: name2, Group: grp2})
	if !strings.Contains(err.Error(), "group and name does not match with lockID") {
		t.Error("Delete error", err)

	}

	//Try create task that already exists
	err = manager.Create(modTask)
	if !strings.Contains(err.Error(), "unable to create, definition already exists") {
		t.Error("Create error", err)

	}

	//Create new task
	err = manager.Create(modTask2)
	if err != nil {
		t.Error("Create error", err)

	}

	lockID, err = manager.Lock(taskdata.GroupNameData{Name: "dummy_AA", Group: "test"})

	err = manager.Delete(lockID, taskdata.GroupNameData{Name: "dummy_AA", Group: "test"})
	if err != nil {
		t.Error("Delete error", err)
	}

}
func TestGetTask(t *testing.T) {

	path, _ := filepath.Abs("../../../def/")
	manager, err := NewManager(path)
	if err != nil {
		t.Fatal("unable to intialize manager")
	}

	tlist := []taskdata.GroupNameData{{Name: "dummy_01", Group: "test"}, {Name: "task_that_does_not_exists", Group: "test"}}

	result := manager.GetTasks(tlist...)
	if len(result) != len(tlist) {
		t.Error("unexpected result, expected :", len(tlist), "got:", len(result))
	}

	if result[0].Result == false || result[1].Result == true {
		t.Error("unexpected values expected:", true, false, "got:", result[0].Result, result[1].Result)
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

func TestMarshalTests2(t *testing.T) {

	data, _ := WriteDefinitionFile(expect3)
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

	schdata := SchedulingData{FromTime: "10:30", ToTime: "11:30"}

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_time_span", "description").
		WithSchedule(schdata).
		WithFlags([]FlagData{{Name: "FLAG01", Type: FlagShared}}).
		WithConfirm().WithRetention(1).WithVariables([]VariableData{{Name: "%%var", Value: "xx"}}).Build()

	if err != nil {
		t.Error("task builder error")
	}

	from, to := def.TimeSpan()
	if from.String() != "10:30" || to.String() != "11:30" {
		t.Error("Unexpected values:", from, to)
	}
}

func TestGetAction(t *testing.T) {

	schdata := SchedulingData{FromTime: "10:30", ToTime: "11:30"}

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_time_span", "description").
		WithSchedule(schdata).
		WithFlags([]FlagData{{Name: "FLAG01", Type: FlagShared}}).
		WithConfirm().WithRetention(1).WithVariables([]VariableData{{Name: "%%var", Value: "xx"}}).Build()

	if err != nil {
		t.Error("task builder error")
	}

	if def.Action() != "" {
		t.Error("Unexpected value")
	}
}

func TestExpandVariable(t *testing.T) {

	expect := "OVS_VARIBALE"
	variable := VariableData{Name: "%%VARIBALE", Value: ""}
	if variable.Expand() != expect {
		t.Error("Unexpected value expected:", expect, "actual", variable.Expand())
	}

}

func TestBuilder(t *testing.T) {

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_04", "description").
		WithFlags([]FlagData{{Name: "FLAG01", Type: FlagShared}}).
		WithConfirm().WithRetention(1).WithVariables([]VariableData{{Name: "%%var", Value: "xx"}}).Build()

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
			Values:    []string{"1", "3", "5"},
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
			Values:    []string{"1", "3", "5"},
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
			Values:    []string{"1", "3", "5"},
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
			Values:    []string{"1", "3", "5"},
		},
	}
	return
}

package taskdef

import (
	"os"
	"testing"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/overseer/internal/unique"
	"github.com/przebro/overseer/overseer/taskdata"
)

func TestNewManager_Errors(t *testing.T) {

	_, err := NewManager("path_thath_does_not_exists", logger.NewTestLogger())
	if err == nil {
		t.Error("unexpected result")
	}

	workdir, _ := os.Getwd()
	_, err = NewManager(workdir, logger.NewTestLogger())
	if err == nil {
		t.Error("unexpected result")
	}

}

func TestCreateManager(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())
	if manager == nil {
		t.Error("unexpected result")
	}
}

func TestWriteActiveDefinition(t *testing.T) {

	m, err := NewManager(managerPath, logger.NewTestLogger())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	builder := DummyTaskBuilder{}
	task, err := builder.WithBase("test", "dummy_0A", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()
	if err != nil {
		t.Error("unexpected result:", err)
	}

	err = m.WriteActiveDefinition(task, unique.NewID())
	if err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestRemoveActivDefinition(t *testing.T) {

	m, err := NewManager(managerPath, logger.NewTestLogger())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	builder := DummyTaskBuilder{}
	task, err := builder.WithBase("test", "dummy_0A", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()
	if err != nil {
		t.Error("unexpected result:", err)
	}

	id := unique.NewID()

	err = m.WriteActiveDefinition(task, id)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	err = m.RemoveActiveDefinition(id.Hex())

	if err != nil {
		t.Error("unexpected result:", err)
	}

}
func TestGetActivDefinition_Errors(t *testing.T) {

	m, err := NewManager(managerPath, logger.NewTestLogger())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if _, err := m.GetActiveDefinition(""); err != ErrReferenceEmpty {
		t.Error("unexpected result:", err, "expected:", ErrReferenceEmpty)
	}

	_, err = m.GetActiveDefinition("id_that_does_not_exists")

	if _, ok := err.(*os.PathError); !ok {
		t.Error("unexpected result:", err, "expected: os.PathError")
	}

}
func TestGetActivDefinition(t *testing.T) {

	m, err := NewManager(managerPath, logger.NewTestLogger())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	builder := DummyTaskBuilder{}
	task, err := builder.WithBase("test", "dummy_0A", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()
	if err != nil {
		t.Error("unexpected result:", err)
	}

	id := unique.NewID()

	err = m.WriteActiveDefinition(task, id)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = m.GetActiveDefinition(id.Hex())

	if err != nil {
		t.Error("unexpected result:", err)
	}

}

func TestManager_Update_Error_EmptyName(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	//builder := &DummyTaskBuilder{}

	emptyDef := &baseTaskDefinition{}

	if err := manager.Update(emptyDef); err != ErrTaskNameEmpty {
		t.Error("unexpected result:", err, "expected:", ErrTaskNameEmpty)
	}
}

func TestManager_Update_Error_EmptyRev(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	emptyDef := &baseTaskDefinition{Name: "default_name"}

	if err := manager.Update(emptyDef); err != ErrTaskRevEmpty {
		t.Error("unexpected result:", err, "expected:", ErrTaskRevEmpty)
	}
}

func TestManager_Update_Error_InvalidRev(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	emptyDef := &baseTaskDefinition{Name: "default_name", Revision: "abc"}

	if err := manager.Update(emptyDef); err != ErrTaskRevInvalid {
		t.Error("unexpected result:", err, "expected:", ErrTaskRevInvalid)
	}
}

func TestManager_Update_Error_RevDiff(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	builder := &DummyTaskBuilder{}
	def, _ := builder.WithBase("test", "dummy_update_01", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()

	manager.Create(def)

	invalidDef, _ := def.(*baseTaskDefinition)
	invalidDef.Revision = "dummy_update_01@test@123456"

	if err := manager.Update(invalidDef); err != ErrTaskRevDiff {
		t.Error("unexpected result:", err, "expected:", ErrTaskRevDiff)
	}

	manager.Delete(taskdata.GroupNameData{Name: "dummy_update_01", GroupData: taskdata.GroupData{Group: "test"}})
}

func TestManager_Update_Error_AlreadyExists(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	builder := &DummyTaskBuilder{}
	def, _ := builder.WithBase("test", "dummy_update_01", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()
	def2, _ := builder.WithBase("test", "dummy_update_02", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()

	manager.Create(def)
	manager.Create(def2)

	invalidDef, _ := def.(*baseTaskDefinition)
	invalidDef.Name = "dummy_update_02"

	if err := manager.Update(invalidDef); err != ErrTaskRename {
		t.Error("unexpected result:", err, "expected:", ErrTaskRename)
	}

	manager.Delete(taskdata.GroupNameData{Name: "dummy_update_01", GroupData: taskdata.GroupData{Group: "test"}})
	manager.Delete(taskdata.GroupNameData{Name: "dummy_update_02", GroupData: taskdata.GroupData{Group: "test"}})
}

func TestManager_Update(t *testing.T) {

	var ok bool
	var base *baseTaskDefinition

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	builder := &DummyTaskBuilder{}

	def, err := builder.WithBase("test", "dummy_update_03", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if err := manager.Create(def); err != nil {
		t.Error("unexpected result:", err)
	}

	if base, ok = def.(*baseTaskDefinition); !ok {
		t.Error("unexpected result")
	}

	base.Cyclics = CyclicTaskData{IsCycle: true, MaxRuns: 10, RunFrom: CycleFromEnd, TimeInterval: 10}

	err = manager.Update(def)
	if err != nil {
		t.Error("unexpected result:")
	}

	manager.Delete(taskdata.GroupNameData{Name: "dummy_update_03", GroupData: taskdata.GroupData{Group: "test"}})
}

func TestManager_GetGroups_Errors(t *testing.T) {
	manager := &taskManager{dirPath: "path_that_does_not_exists", log: logger.NewTestLogger()}
	_, err := manager.GetGroups()
	if err != ErrGroupDirInvalid {
		t.Error("unexpected result:", err, "expected:", ErrGroupDirInvalid)
	}
}

func TestManager_DeleteGroup_Errors(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	err := manager.DeleteGroup("")
	if err.Error() != "group name cannot be empty" {
		t.Error("Unexpected error")
	}
}

func TestManager_DeleteGroup_ErrorNotEmpty(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	if err := manager.DeleteGroup("test"); err != ErrGroupDirNotEmpty {
		t.Error("unexpected result:", err, "expected:", ErrGroupDirNotEmpty)
	}
}

func TestManager_DeleteGroup_ErrorGroupNotExists(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	if err := manager.DeleteGroup("test2"); err != ErrGroupNotExists {
		t.Error("unexpected result:", err, "expected:", ErrGroupNotExists)
	}
}

func TestManagerGroups(t *testing.T) {

	manager, _ := NewManager(managerPath, logger.NewTestLogger())

	groups, _ := manager.GetGroups()

	if len(groups) != groupsDircetories {
		t.Error("invalid group names")
	}

	err := manager.CreateGroup("test2")
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

	//cleanup
	manager.DeleteGroup("test2")

}

func Test_getNameGroupIdFromDefinition_ErrorEmpty(t *testing.T) {

	def := &baseTaskDefinition{}

	if _, _, _, err := getNameGroupIdFromDefinition(def); err != ErrTaskRevEmpty {
		t.Error("unexpected result:", err, "expected:", ErrTaskRevEmpty)
	}
}

func Test_getNameGroupIdFromDefinition_ErrorInvalid(t *testing.T) {

	def := &baseTaskDefinition{Revision: "abc"}

	if _, _, _, err := getNameGroupIdFromDefinition(def); err != ErrTaskRevInvalid {
		t.Error("unexpected result:", err, "expected:", ErrTaskRevInvalid)
	}
}

func Test_getNameGroupIdFromDefinition(t *testing.T) {

	builder := DummyTaskBuilder{}
	def, err := builder.WithBase("test", "dummy_0A", "").WithSchedule(SchedulingData{OrderType: OrderingManual}).Build()
	if err != nil {
		t.Error("unexpected result:", err)
	}

	name, group, rev, _ := getNameGroupIdFromDefinition(def)

	actual := name + "@" + group + "@" + rev

	if actual != def.Rev() {
		t.Error("unexpected result:", actual, "expected:", def.Rev())
	}

}

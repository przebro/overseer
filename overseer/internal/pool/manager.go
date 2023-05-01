package pool

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/internal/journal"
	"github.com/przebro/overseer/overseer/internal/pool/activetask"
	"github.com/przebro/overseer/overseer/internal/pool/models"
	"github.com/przebro/overseer/overseer/internal/pool/states"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/taskdata"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	//ErrUnableFindDef - a definition was not found
	ErrUnableFindDef = errors.New("unable to find definition")
	//ErrInvalidStatus - invalid status
	ErrInvalidStatus = errors.New("invalid status")
	//ErrUnableFindGroup - a group was not found
	ErrUnableFindGroup = errors.New("unable to find group")
	//ErrUnableFindTask 0 task was not found
	ErrUnableFindTask = errors.New("unable to find task")
)

// ActiveTaskPoolManager - Manages tasks in Active task pool
type ActiveTaskPoolManager struct {
	log      zerolog.Logger
	tdm      taskdef.TaskDefinitionManager
	pool     *ActiveTaskPool
	sequence SequenceGenerator
}

// NewActiveTaskPoolManager - Creates a new Managager
func NewActiveTaskPoolManager(
	tdm taskdef.TaskDefinitionManager,
	pool *ActiveTaskPool,
	provider *datastore.Provider) (*ActiveTaskPoolManager, error) {

	lg := log.With().Str("component", "pool-manager").Logger()

	var err error
	manager := &ActiveTaskPoolManager{}
	manager.tdm = tdm
	manager.pool = pool
	manager.log = lg

	if manager.sequence, err = NewSequenceGenerator(provider); err != nil {
		lg.Err(err).Msg("unable to create sequence generator")
		return nil, err
	}

	return manager, nil
}

// Order - Orders a new task, this method checks if all precoditions are met before it adds a new task
func (manager *ActiveTaskPoolManager) Order(task taskdata.GroupNameData, odate date.Odate, username string) (string, error) {

	var err error
	var definition *taskdef.TaskDefinition
	if definition, err = manager.tdm.GetTask(task); err != nil {
		return "", ErrUnableFindDef
	}

	orderID, descr := manager.orderDefinition(definition, odate, false, username)
	if descr != "" {
		return descr, errors.New("failed order task")
	}

	return string(orderID), nil
}

// OrderGroup - Orders all tasks from group, this method checks if all precoditions are met before it adds a new task
func (manager *ActiveTaskPoolManager) OrderGroup(groupdata taskdata.GroupData, odate date.Odate, username string) ([]string, error) {

	var err error
	var result = []string{}
	var grps []taskdata.GroupNameData
	var definition *taskdef.TaskDefinition

	if grps, err = manager.tdm.GetTasksFromGroup([]string{groupdata.Group}); err != nil {
		return []string{}, ErrUnableFindGroup
	}

	for _, d := range grps {

		if definition, err = manager.tdm.GetTask(d); err != nil {
			manager.log.Error().Err(err).Str("group", d.Group).Str("task", d.Name).Msg("get task definition")
			continue
		}

		orderID, descr := manager.orderDefinition(definition, odate, false, username)
		if descr != "" {
			result = append(result, descr)
			continue
		}
		result = append(result, string(orderID))
	}

	return result, nil
}

// ForceGroup - forcefully orders all tasks from group, this method checks if all precoditions are met before it adds a new task
func (manager *ActiveTaskPoolManager) ForceGroup(groupdata taskdata.GroupData, odate date.Odate, username string) ([]string, error) {

	var err error
	var result = []string{}
	var grps []taskdata.GroupNameData
	var definition *taskdef.TaskDefinition

	if grps, err = manager.tdm.GetTasksFromGroup([]string{groupdata.Group}); err != nil {
		return []string{}, ErrUnableFindGroup
	}

	for _, d := range grps {

		definition, _ = manager.tdm.GetTask(d)

		orderID, descr := manager.orderDefinition(definition, odate, true, username)
		if descr != "" {
			result = append(result, descr)
			continue
		}
		result = append(result, string(orderID))
	}

	return result, nil
}

// Force - forcefully orders a new task, this method does not check for precondtions
func (manager *ActiveTaskPoolManager) Force(task taskdata.GroupNameData, odate date.Odate, username string) (string, error) {

	var err error
	var definition *taskdef.TaskDefinition
	if definition, err = manager.tdm.GetTask(task); err != nil {
		return "", ErrUnableFindDef
	}

	orderID, descr := manager.orderDefinition(definition, odate, true, username)
	if descr != "" {
		return descr, errors.New("failed force task")
	}

	return string(orderID), nil
}

// Rerun - Orders task again
func (manager *ActiveTaskPoolManager) Rerun(id unique.TaskOrderID, username string) (string, error) {

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), ErrUnableFindTask
	}

	state := task.State()
	if state == models.TaskStateEndedNotOk || state == models.TaskStateEndedOk {
		task.SetExecutionID()
	} else {
		return fmt.Sprintf("rerun task:%s failed, invalid status", id), ErrInvalidStatus
	}

	manager.pool.journal.PushJournalMessage(task.OrderID(), task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskRerun, username))
	return fmt.Sprintf("rerun task:%s ok", id), nil

}

// Enforce - enforces task execution
func (manager *ActiveTaskPoolManager) Enforce(id unique.TaskOrderID, username string) (string, error) {

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), ErrUnableFindTask
	}

	state := task.State()
	if state != models.TaskStateWaiting {
		return fmt.Sprintf("enforce task:%s failed, invalid status", id), ErrInvalidStatus
	}

	manager.pool.enforceTask(id)

	manager.pool.journal.PushJournalMessage(task.OrderID(), task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskEnforce, username))

	return fmt.Sprintf("enforce task:%s ok", id), nil
}

// SetOk - Sets task to EndedOk status
func (manager *ActiveTaskPoolManager) SetOk(id unique.TaskOrderID, username string) (string, error) {

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), ErrUnableFindTask
	}

	state := task.State()
	if state == models.TaskStateEndedNotOk {
		task.SetState(models.TaskStateEndedOk)
	} else {
		return fmt.Sprintf("set to OK task:%s failed, invalid status", id), ErrInvalidStatus
	}

	manager.pool.journal.PushJournalMessage(task.OrderID(), task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskSetOK, username))
	return fmt.Sprintf("set to ok task:%s ok", id), nil
}

// Hold - Holds the taskdef. It will be not processed durning cycle
func (manager *ActiveTaskPoolManager) Hold(id unique.TaskOrderID, username string) (string, error) {

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), ErrUnableFindTask
	}

	if task.IsHeld() {
		return fmt.Sprintf("hold task:%s failed, invalid status", id), ErrInvalidStatus

	}
	task.Hold()
	manager.pool.journal.PushJournalMessage(task.OrderID(), task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskHeld, username))

	return fmt.Sprintf("hold task:%s ok", id), nil
}

// Free - Frees a holded task
func (manager *ActiveTaskPoolManager) Free(id unique.TaskOrderID, username string) (string, error) {

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), ErrUnableFindTask
	}

	if !task.IsHeld() {
		return fmt.Sprintf("free task:%s failed, task is not held", id), ErrInvalidStatus
	}

	task.Free()
	manager.pool.journal.PushJournalMessage(task.OrderID(), task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskFreed, username))

	return fmt.Sprintf("free task:%s ok", id), nil
}

// Confirm - Manually Confirms a task
func (manager *ActiveTaskPoolManager) Confirm(id unique.TaskOrderID, username string) (string, error) {

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), ErrUnableFindTask
	}

	result := task.SetConfirm()
	if !result {
		return fmt.Sprintf("task with id:%s already confirmed", id), ErrInvalidStatus
	}

	manager.pool.journal.PushJournalMessage(task.OrderID(), task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskConfirmed, username))

	return fmt.Sprintf("confirm task:%s ok", id), nil
}

func (manager *ActiveTaskPoolManager) OrderNewTasks() int {

	ordered := 0

	manager.log.Info().Msg("Ordering new tasks")

	groups, _ := manager.tdm.GetGroups()

	taskData, _ := manager.tdm.GetTasksFromGroup(groups)

	result := manager.tdm.GetTasks(taskData...)

	for _, t := range result {

		//It is a new day procedure so skip tasks that are ordered manually
		if t.Schedule.OrderType != taskdef.OrderingManual {
			manager.orderDefinition(&t, date.CurrentOdate(), false, "daily procedure")
			ordered++
		}
	}
	return ordered
}

// orderDefinition - Adds a new task to the Active Task Pool
// this method performs all checks
func (manager *ActiveTaskPoolManager) orderDefinition(def *taskdef.TaskDefinition, odate date.Odate, force bool, username string) (unique.TaskOrderID, string) {

	manager.log.Info().Str("", "").Str("odate", string(odate)).Msg("order definition")

	ctx := states.TaskOrderContext{
		Def:              def,
		IgnoreCalendar:   force,
		IgnoreSubmission: force,
		Odate:            odate,
		State:            &states.OstateCheckOtype{},
		CurrentOdate:     date.CurrentOdate(),
		Reason:           make([]string, 0),
		Log:              manager.log,
	}
	n, g, _ := def.GetInfo()

	for ctx.State.ProcessState(&ctx) {
	}

	if !ctx.IsSubmited {
		return "", strings.Join(ctx.Reason, ",")

	}

	refID := unique.NewID()

	if err := manager.tdm.WriteActiveDefinition(def, refID); err != nil {
		manager.log.Error().Err(err).Msg("push definition to pool failed")
	}

	orderID := manager.sequence.Next()
	task := activetask.NewActiveTask(orderID, odate, def, refID)
	manager.pool.addTask(orderID, task)

	if force {
		manager.log.Info().Msg(fmt.Sprintf("Task %s from gorup %s forced with id:%s odate:%s", n, g, orderID, odate))
		manager.pool.journal.PushJournalMessage(orderID, task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskForced, username, odate))
	} else {
		manager.log.Info().Msg(fmt.Sprintf("Task %s from gorup %s ordered with id:%s odate:%s", n, g, orderID, odate))
		manager.pool.journal.PushJournalMessage(orderID, task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskOrdered, username, odate))
	}

	return orderID, strings.Join(ctx.Reason, ",")
}

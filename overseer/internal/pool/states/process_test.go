package states

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/pool/activetask"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockSequence struct {
	seq unique.TaskOrderID
}

func (m *mockSequence) Next() unique.TaskOrderID {

	return m.seq
}

type mockResourceManager struct {
	mock.Mock
}

func (m *mockResourceManager) ProcessReleaseFlag(input []string) (bool, []string) {
	return true, []string{}
}
func (m *mockResourceManager) ProcessAcquireFlag(input []types.FlagModel) (bool, []string) {
	return true, []string{}
}
func (m *mockResourceManager) ProcessTicketAction(tickets []types.TicketActionModel) bool {
	return true
}
func (m *mockResourceManager) CheckTickets(in []types.CollectedTicketModel) []types.CollectedTicketModel {
	args := m.Called(in)
	return args.Get(0).([]types.CollectedTicketModel)
}

type mockJournalWriter struct {
	mock.Mock
}

func (m *mockJournalWriter) PushJournalMessage(ID unique.TaskOrderID, execID string, t time.Time, msg string) {

}

type TaskStateTestSuite struct {
	suite.Suite
	mRmanager *mockResourceManager
	seq       *mockSequence
	builder   taskdef.TaskBuilder
	now       time.Time
	odate     date.Odate
	Context   func() TaskExecutionContext
}

func TestTaskStateTestSuite(t *testing.T) {
	suite.Run(t, new(TaskStateTestSuite))
}
func (suite *TaskStateTestSuite) SetupSuite() {
	suite.seq = &mockSequence{seq: "12345"}
	suite.mRmanager = &mockResourceManager{}
	suite.builder = taskdef.NewBuilder()
	suite.now = time.Now()
	suite.odate = date.CurrentOdate()

	suite.Context = func() TaskExecutionContext {
		return TaskExecutionContext{
			Log:      log.Logger,
			Odate:    suite.odate,
			Time:     suite.now,
			Rmanager: suite.mRmanager,
			IsInTime: true,
			Journal:  &mockJournalWriter{},
		}
	}

	suite.mRmanager.On("CheckTickets", []types.CollectedTicketModel{
		{Name: "AND_01_TESTABC01", Odate: date.CurrentOdate(), Exists: false},
		{Name: "AND_01_TESTABC02", Odate: date.CurrentOdate(), Exists: false},
	}).Return(
		[]types.CollectedTicketModel{
			{Name: "AND_01_TESTABC01", Odate: date.CurrentOdate(), Exists: false},
			{Name: "AND_01_TESTABC02", Odate: date.CurrentOdate(), Exists: false},
		},
	)
	suite.mRmanager.On("CheckTickets", []types.CollectedTicketModel{
		{Name: "AND_02_TESTABC01", Odate: date.CurrentOdate(), Exists: false},
		{Name: "AND_02_TESTABC02", Odate: date.CurrentOdate(), Exists: false},
	}).Return(
		[]types.CollectedTicketModel{
			{Name: "AND_02_TESTABC01", Odate: date.CurrentOdate(), Exists: true},
			{Name: "AND_02_TESTABC02", Odate: date.CurrentOdate(), Exists: true},
		},
	)
	suite.mRmanager.On("CheckTickets", []types.CollectedTicketModel{
		{Name: "OR_03_TESTABC01", Odate: date.CurrentOdate(), Exists: false},
		{Name: "OR_03_TESTABC02", Odate: date.CurrentOdate(), Exists: false},
	}).Return(
		[]types.CollectedTicketModel{
			{Name: "OR_03_TESTABC01", Odate: date.CurrentOdate(), Exists: true},
			{Name: "OR_03_TESTABC02", Odate: date.CurrentOdate(), Exists: false},
		},
	)

	suite.mRmanager.On("CheckTickets", []types.CollectedTicketModel{
		{Name: "EX_04_TESTABC04A", Odate: date.CurrentOdate(), Exists: false},
		{Name: "EX_04_TESTABC04B", Odate: date.CurrentOdate(), Exists: false},
	}).Return(
		[]types.CollectedTicketModel{
			{Name: "EX_04_TESTABC04A", Odate: date.CurrentOdate(), Exists: true},
			{Name: "EX_04_TESTABC04B", Odate: date.CurrentOdate(), Exists: false},
		},
	)

	suite.mRmanager.On("CheckTickets", []types.CollectedTicketModel{
		{Name: "EX_05_TESTABC05", Odate: date.CurrentOdate(), Exists: false},
		{Name: "EX_05_TESTABC05", Odate: date.CurrentOdate(), Exists: false},
	}).Return(
		[]types.CollectedTicketModel{
			{Name: "EX-05-TESTABC04", Odate: date.CurrentOdate(), Exists: true},
			{Name: "EX-05-TESTABC04", Odate: date.CurrentOdate(), Exists: false},
		},
	)
}

func (suite *TaskStateTestSuite) TestState_CheckTime() {

	type testCase struct {
		taskdef.SchedulingData
		expected   bool
		isEnforced bool
	}

	now := time.Now()
	nowPlus2 := now.Add(2 * time.Minute)
	nowPlus4 := now.Add(4 * time.Minute)
	nowMinus2 := now.Add(-2 * time.Minute)
	nowMinus4 := now.Add(-4 * time.Minute)

	strn := types.HourMinTime(fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()))
	strnp10 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowPlus2.Hour(), nowPlus2.Minute()))
	strnp20 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowPlus4.Hour(), nowPlus4.Minute()))
	strnm10 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowMinus2.Hour(), nowMinus2.Minute()))
	strnm20 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowMinus4.Hour(), nowMinus4.Minute()))

	testCases := []testCase{
		{taskdef.SchedulingData{FromTime: "", ToTime: "", OrderType: taskdef.OrderingManual}, true, false},
		{taskdef.SchedulingData{FromTime: strnp10, ToTime: "", OrderType: taskdef.OrderingDaily}, false, false}, // now -> from "-" -> to "-"
		{taskdef.SchedulingData{FromTime: "", ToTime: strnp20, OrderType: taskdef.OrderingDaily}, true, false},
		{taskdef.SchedulingData{FromTime: strnm10, ToTime: strnp10, OrderType: taskdef.OrderingDaily}, true, false},  //from "-" -> now -> to "-"
		{taskdef.SchedulingData{FromTime: strnm20, ToTime: strnm10, OrderType: taskdef.OrderingDaily}, false, false}, //from "-" -> to "-" -> now
		{taskdef.SchedulingData{FromTime: "", ToTime: strnm10, OrderType: taskdef.OrderingDaily}, false, false},
		{taskdef.SchedulingData{FromTime: strn, ToTime: "", OrderType: taskdef.OrderingDaily}, true, false},
		{taskdef.SchedulingData{FromTime: "", ToTime: strn, OrderType: taskdef.OrderingDaily}, false, false},
		{taskdef.SchedulingData{FromTime: strnp10, ToTime: strnp20, OrderType: taskdef.OrderingDaily}, true, true}, //from '-' -> to now
	}

	builder := taskdef.NewBuilder()

	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{FromTime: "", ToTime: "", OrderType: taskdef.OrderingManual}).
		Build()

	if err != nil {
		suite.Fail("Unable to construct task")
	}

	state := OstateCheckTime{}

	ctx := suite.Context()

	for i, caseT := range testCases {

		definition, err = builder.FromTemplate(definition).WithSchedule(caseT.SchedulingData).Build()

		suite.Nil(err)

		ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())
		ctx.IsEnforced = caseT.isEnforced

		state.ProcessState(&ctx)
		suite.Equal(caseT.expected, ctx.IsInTime, fmt.Sprint("num:", i, ",expected result:", caseT.expected, "actual:", ctx.IsInTime, " ", "now:", strn, " ", caseT.SchedulingData))

	}
}

func (suite *TaskStateTestSuite) TestState_CheckCond_AndRelationNegativeResult() {

	var result bool

	builder := taskdef.NewBuilder()
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithInTicekts([]taskdef.InTicketData{
			{Name: "AND_01_TESTABC01", Odate: date.OdateValueDate},
			{Name: "AND_01_TESTABC02", Odate: date.OdateValueDate},
		}, "AND").
		Build()

	suite.Nil(err)

	ctx := suite.Context()

	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := OstateCheckConditions{}
	result = state.ProcessState(&ctx)

	suite.Equal(false, result)
}
func (suite *TaskStateTestSuite) TestState_CheckCond_AndRelationPositiveResult() {

	var result bool

	builder := taskdef.NewBuilder()
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithInTicekts([]taskdef.InTicketData{
			{Name: "AND_02_TESTABC01", Odate: date.OdateValueDate},
			{Name: "AND_02_TESTABC02", Odate: date.OdateValueDate},
		}, "AND").
		Build()

	suite.Nil(err)

	ctx := suite.Context()
	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := OstateCheckConditions{}
	result = state.ProcessState(&ctx)

	suite.Equal(true, result)
}
func (suite *TaskStateTestSuite) TestState_CheckCond_OrRelationPositiveResult() {
	var result bool

	definition, err := suite.builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithInTicekts([]taskdef.InTicketData{
			{Name: "OR_03_TESTABC01", Odate: date.OdateValueDate},
			{Name: "OR_03_TESTABC02", Odate: date.OdateValueDate},
		}, "OR").
		Build()

	suite.Nil(err)

	ctx := suite.Context()

	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := OstateCheckConditions{}
	result = state.ProcessState(&ctx)

	suite.Equal(true, result)
}
func (suite *TaskStateTestSuite) TestState_CheckCond_IsEnforced() {
	var result bool

	builder := taskdef.NewBuilder()
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).Build()

	suite.Nil(err)

	ctx := suite.Context()
	ctx.IsEnforced = true
	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := OstateCheckConditions{}
	result = state.ProcessState(&ctx)

	suite.Equal(true, result)
}
func (suite *TaskStateTestSuite) TestState_CheckCondOr_Expression() {

	var result bool

	definition, err := suite.builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithInTicekts([]taskdef.InTicketData{
			{Name: "EX_04_TESTABC04A", Odate: date.OdateValueDate},
			{Name: "EX_04_TESTABC04B", Odate: date.OdateValueDate},
		}, "EX_04_TESTABC04A.ODATE || EX_04_TESTABC04B.ODATE").
		Build()

	suite.Nil(err)

	ctx := suite.Context()
	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := OstateCheckConditions{}
	result = state.ProcessState(&ctx)

	suite.Equal(true, result)

}

func (suite *TaskStateTestSuite) TestState_CheckCondAnd_Expression() {

	var result bool

	definition, err := suite.builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithInTicekts([]taskdef.InTicketData{
			{
				Name: "EX_05_TESTABC05", Odate: date.OdateValueDate,
			},
			{
				Name: "EX_05_TESTABC05", Odate: date.OdateValueDate,
			},
		}, "EX_05_TESTABC05.ODATE && EX_05_TESTABC05.ODATE").
		Build()

	suite.Nil(err)

	ctx := suite.Context()
	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := OstateCheckConditions{}
	result = state.ProcessState(&ctx)

	suite.Equal(false, result)

}

func (suite *TaskStateTestSuite) TestState_Confirmed() {
	var result bool

	definition, err := suite.builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithConfirm().Build()
	suite.Nil(err)

	ctx := suite.Context()
	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	state := OstateConfirm{}
	result = state.ProcessState(&ctx)

	suite.Equal(false, result)

	definition, err = suite.builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		Build()
	suite.Nil(err)

	ctx.Task = activetask.NewActiveTask(suite.seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	state = OstateConfirm{}
	result = state.ProcessState(&ctx)

	suite.Equal(true, result)
}

/*
func TestNewCyclicTask(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: 1,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	cycle := task.CycleData()

	if !cycle.IsCyclic {
		t.Error("unexpected value:", cycle.IsCyclic, "expected:", false)
	}
	if cycle.MaxRun != 3 {
		t.Error("unexpected value:", cycle.MaxRun, "expected:", 3)
	}

	if cycle.RunInterval != 1 {
		t.Error("unexpected value:", cycle.RunInterval, "expected:", 1)
	}

	if cycle.RunFrom != string(taskdef.CycleFromStart) {
		t.Error("unexpected value:", cycle.RunFrom, "expected:", taskdef.CycleFromStart)
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = true
	state := ostatePostProcessing{}
	result := state.processState(&ctx)
	fmt.Println(result)

}

func TestCyclicState_Enforced(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: 1,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = true
	state := ostateCheckCyclic{}
	result := state.processState(&ctx)
	if !result {
		t.Error("unexpected result:", result, "expected:", true)
	}

}

func TestCyclicState_Normal(t *testing.T) {

		builder := taskdef.DummyTaskBuilder{}
		definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
			WithSchedule(taskdef.SchedulingData{
				OrderType: taskdef.OrderingDaily,
			}).Build()

		if err != nil {
			t.Log(err)
		}

		task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

		ctx := TaskExecutionContext{
			log:        log,
			odate:      date.CurrentOdate(),
			dispatcher: mDispatcher,
			time:       time.Now(),
		}

		ctx.task = task
		ctx.isEnforced = false
		state := ostateCheckCyclic{}
		result := state.processState(&ctx)
		if !result {
			t.Error("unexpected result:", result, "expected:", true)
		}
	}

func TestCyclicState_BeforeNextRun(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: 1,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	now := time.Now().Add(2 * time.Minute)
	task.cycle.NextRun = types.FromTime(now)

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = false
	state := ostateCheckCyclic{}
	result := state.processState(&ctx)
	if result {
		t.Error("unexpected result:", result, "expected:", false)
	}

}

func TestStatesExecEndHold(t *testing.T) {

	var result bool

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).
		WithOutTickets([]taskdef.OutTicketData{{Action: "REM", Name: "TEST", Odate: "ODATE"}}).
		Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := &ostateStarting{}
	runState := &ostateExecuting{}
	postState := &ostatePostProcessing{}

	result = state.processState(&ctx)

	if result != true && ctx.state != runState {
		t.Error("expected result: ", true, " actual:", result, " ", ctx.state)
	}

	mDispatcher.processNotEnded = true
	mDispatcher.withError = false
	result = runState.processState(&ctx)

	if result != false && ctx.state != runState {
		t.Error("expected result: ", true, " actual:", result, " ", ctx.state)
	}

	mDispatcher.withError = true
	result = runState.processState(&ctx)

	if result != false && ctx.state != runState {
		t.Error("expected result: ", false, " actual:", result, " ", ctx.state)
	}

	mDispatcher.withError = false
	mDispatcher.processNotEnded = false
	result = postState.processState(&ctx)
	if result != false {
		t.Error("expected result: ", false, " actual:", result)
	}
	ctx.task.SetState(TaskStateEndedOk)

	if ctx.task.State() != TaskStateEndedOk {
		t.Error("expected result: ", TaskStateEndedOk, " actual:", ctx.task.State())
	}

	if ctx.task.OrderDate() != date.CurrentOdate() {
		t.Error("expected result: ", date.CurrentOdate(), " actual:", ctx.task.OrderDate())
	}

	if ctx.task.RunNumber() != 1 {
		t.Error("expected result: ", 1, " actual:", ctx.task.RunNumber())
	}

	ctx.task.Hold()
	if ctx.task.IsHeld() != true {
		t.Error("expected result: ", true, " actual:", ctx.task.IsHeld())
	}

	ctx.task.Free()
	if ctx.task.IsHeld() != false {
		t.Error("expected result: ", false, " actual:", ctx.task.IsHeld())
	}

}

func TestGetProcState(t *testing.T) {

	stateWaiting := &ostateConfirm{}
	stateExecute := &ostateExecuting{}

	res := getProcessState(TaskStateWaiting, false)
	if res == nil && res != stateWaiting {
		t.Error("expected result: ", stateWaiting, " actual:", res)
	}

	res = getProcessState(TaskStateExecuting, false)
	if res == nil && res != stateExecute {
		t.Error("expected result: ", stateExecute, " actual:", res)
	}

	res = getProcessState(TaskStateStarting, false)
	if res != nil {
		t.Error("expected result: ", nil, " actual:", res)
	}

	res = getProcessState(TaskStateWaiting, true)
	if res != nil {
		t.Error("expected result: ", nil, " actual:", res)
	}

}

func TestStrTimeToInt(t *testing.T) {

		h, m := strTimeToInt("12:23")
		if h != 12 && m != 23 {
			t.Error("unexpected result")
		}
	}

func TestCalcRealOdate(t *testing.T) {

	current := date.Odate("20201101")
	otype := []taskdef.SchedulingOption{
		taskdef.OrderingManual,
		taskdef.OrderingDaily,
		taskdef.OrderingWeek,
		taskdef.OrderingDayOfMonth,
		taskdef.OrderingExact,
		taskdef.OrderingFromEnd,
	}
	schedule := taskdef.SchedulingData{}
	for x := range otype {
		schedule.OrderType = otype[x]
		//ODATE means same date as current
		result := calcRealOdate(current, date.OdateValueDate, schedule)
		if result != current {
			t.Error("unexpected error, expecting odate:", current, "actual:", result)
		}
	}

	current = date.Odate("20201101")
	schedule.OrderType = taskdef.OrderingManual
	//Tasks ordered manually don't depend on calendar so PREV or NEXT means -1 or +1
	result := calcRealOdate(current, date.OdateValueAny, schedule)
	if result != "" {
		t.Error("unexpected result:", result, "expected empty")
	}

	current = date.Odate("20201101")
	schedule.OrderType = taskdef.OrderingManual
	//Tasks ordered manually don't depend on calendar so PREV or NEXT means -1 or +1
	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20201031" {
		t.Error("unexpected result:", result, "expected 20201001")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != "20201102" {
		t.Error("unexpected result:", result, "expected 20201102")
	}

	//For exact date with only a single value there NEXT and PREV should resolve to +1 and +1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-11-01"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != "20201102" {
		t.Error("unexpected result:", result, "expected 20201102")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20201031" {
		t.Error("unexpected result:", result, "expected 20201030")
	}

	//For exact date with current date first,  NEXT should be computed to defined value and PREV to -1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-11-01", "2020-11-15"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20201115") {
		t.Error("unexpected result:", result, "expected empty odate")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20201031" {
		t.Error("unexpected result:", result, "expected 20201031")
	}

	//For exact date with current date first,  PREV should be computed to defined value and NEXT to -1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-10-15", "2020-11-01"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != "20201102" {
		t.Error("unexpected result:", result, "expected 20201102")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201015") {
		t.Error("unexpected result:", result, "expected empty odate")
	}

	//For exact date with current date in middle of values both PREV and NEXT should be Computed
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-10-15", "2020-11-01", "2020-11-15"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20201115") {
		t.Error("unexpected result:", result, "expected 20201115")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201015") {
		t.Error("unexpected result:", result, "expected 20201015")
	}

	//For a task that was forced on a non scheduled day NEXT and PERV should be computed to +1,-1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-10-15", "2020-11-01"}
	current = date.Odate("20201022")

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20201023") {
		t.Error("unexpected result:", result, "expected 20201101")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201021") {
		t.Error("unexpected result:", result, "expected 20201015")
	}

	current = date.Odate("20200923")

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20200924") {
		t.Error("unexpected result:", result, "expected 20201015")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20200922" {
		t.Error("unexpected result:", result, "expected 20200922")
	}

}

func TestCalcRealOdateFromEnd(t *testing.T) {

		schedule := taskdef.SchedulingData{}

		schedule.OrderType = taskdef.OrderingFromEnd
		schedule.Dayvalues = []int{1}
		schedule.Months = []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		current := date.Odate("20201231")

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		result := calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210131") {
			t.Error("unexpected result:", result, "expected empty odate")
		}

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20201130") {
			t.Error("unexpected result:", result, "expected empty odate")
		}

		schedule.OrderType = taskdef.OrderingFromEnd
		schedule.Dayvalues = []int{1}
		schedule.Months = []time.Month{2, 4, 7}

		current = date.Odate("20200430")

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20200731") {
			t.Error("unexpected result:", result, "expected 20200731")
		}

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		//test leap year
		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20200229") {
			t.Error("unexpected result:", result, "expected 20200229")
		}

		current = date.Odate("20210430")

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210731") {
			t.Error("unexpected result:", result, "expected 20210731")
		}

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		//test leap year
		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210228") {
			t.Error("unexpected result:", result, "expected 20210228")
		}

		schedule.OrderType = taskdef.OrderingFromEnd
		schedule.Dayvalues = []int{2}
		schedule.Months = []time.Month{2, 4, 7}
		current = date.Odate("20210228")

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210429") {
			t.Error("unexpected result:", result, "expected 20210731")
		}

		//from end in simple scenario it simply nth day of a next and previous month if they are available
		//test leap year
		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20200730") {
			t.Error("unexpected result:", result, "expected 20210228")
		}
	}

func TestCalcRealOdateWeek(t *testing.T) {

		schedule := taskdef.SchedulingData{}
		schedule.OrderType = taskdef.OrderingWeek
		schedule.Dayvalues = []int{4, 5}
		schedule.Months = []time.Month{3, 7}

		current := date.Odate("20210701")

		result := calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210702") {
			t.Error("unexpected result:", result, "expected 20210702")
		}

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210326") {
			t.Error("unexpected result:", result, "expected 20210326")
		}

		schedule.Dayvalues = []int{5, 7}

		current = date.Odate("20210314")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210312") {
			t.Error("unexpected result:", result, "expected 20210312")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210319") {
			t.Error("unexpected result:", result, "expected 20210319")
		}

		schedule.Dayvalues = []int{5, 7}
		current = date.Odate("20210328")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210326") {
			t.Error("unexpected result:", result, "expected 20210326")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210702") {
			t.Error("unexpected result:", result, "expected 20210702")
		}

		schedule.Dayvalues = []int{5, 7}
		current = date.Odate("20210326")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210321") {
			t.Error("unexpected result:", result, "expected 20210321")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210328") {
			t.Error("unexpected result:", result, "expected 20210328")
		}

		schedule.Dayvalues = []int{5, 7}
		schedule.Months = []time.Month{1, 3}
		current = date.Odate("20210328")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210326") {
			t.Error("unexpected result:", result, "expected 20210326")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20220102") {
			t.Error("unexpected result:", result, "expected 20220102")
		}

		schedule.Dayvalues = []int{2, 7}
		schedule.Months = []time.Month{1, 3}
		current = date.Odate("20210324")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210323") {
			t.Error("unexpected result:", result, "expected 20210323")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210325") {
			t.Error("unexpected result:", result, "expected 20210325")
		}
	}

	func TestCalcRealOdateMonth(t *testing.T) {
		schedule := taskdef.SchedulingData{}
		schedule.OrderType = taskdef.OrderingDayOfMonth
		schedule.Dayvalues = []int{30}
		schedule.Months = []time.Month{1, 2, 3, 7}
		current := date.Odate("20210330")

		result := calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210130") {
			t.Error("unexpected result:", result, "expected 20210130")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210730") {
			t.Error("unexpected result:", result, "expected 20210730")
		}

		//non scheduled day
		current = date.Odate("20210328")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210327") {
			t.Error("unexpected result:", result, "expected 20210327")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210329") {
			t.Error("unexpected result:", result, "expected 20210329")
		}

		schedule.Dayvalues = []int{31}
		schedule.Months = []time.Month{1, 2, 4, 7}
		current = date.Odate("20210731")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210131") {
			t.Error("unexpected result:", result, "expected 20210131")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20220131") {
			t.Error("unexpected result:", result, "expected 20220131")
		}

		schedule.Dayvalues = []int{29, 30, 31}
		schedule.Months = []time.Month{1, 2, 4}
		current = date.Odate("20210131")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210130") {
			t.Error("unexpected result:", result, "expected 20210129")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210429") {
			t.Error("unexpected result:", result, "expected 20210429")
		}

		schedule.Dayvalues = []int{29, 30, 31}
		schedule.Months = []time.Month{1, 2, 4}
		current = date.Odate("20200131")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20200130") {
			t.Error("unexpected result:", result, "expected 20210129")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20200229") {
			t.Error("unexpected result:", result, "expected 20220229")
		}

		schedule.Dayvalues = []int{1, 17, 28, 29, 30, 31}
		schedule.Months = []time.Month{1, 2, 4}
		current = date.Odate("20210228")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210217") {
			t.Error("unexpected result:", result, "expected 20210217")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210401") {
			t.Error("unexpected result:", result, "expected 20210401")
		}

		schedule.Dayvalues = []int{2, 17, 28, 29, 30, 31}
		schedule.Months = []time.Month{1, 4, 11}
		current = date.Odate("20210102")

		result = calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20201130") {
			t.Error("unexpected result:", result, "expected 20201130")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210117") {
			t.Error("unexpected result:", result, "expected 20210217")
		}

}

	func TestCalcRealOdateDaily(t *testing.T) {
		schedule := taskdef.SchedulingData{}
		schedule.OrderType = taskdef.OrderingDaily
		schedule.Dayvalues = []int{30}
		schedule.Months = []time.Month{1, 3, 7}
		current := date.Odate("20210301")

		result := calcRealOdate(current, date.OdateValuePrev, schedule)
		if result != date.Odate("20210131") {
			t.Error("unexpected result:", result, "expected 20210131")
		}

		result = calcRealOdate(current, date.OdateValueNext, schedule)
		if result != date.Odate("20210302") {
			t.Error("unexpected result:", result, "expected 20210302")
		}
	}

func TestCalcRealOdateRelative(t *testing.T) {

		current := date.Odate("20201105")
		otype := []taskdef.SchedulingOption{
			taskdef.OrderingManual,
			taskdef.OrderingDaily,
			taskdef.OrderingWeek,
			taskdef.OrderingDayOfMonth,
			taskdef.OrderingExact,
			taskdef.OrderingFromEnd,
		}
		schedule := taskdef.SchedulingData{}
		for x := range otype {
			schedule.OrderType = otype[x]

			result := calcRealOdate(current, "-001", schedule)
			if result != "20201104" {
				t.Error("unexpected error, expecting odate:", "20201104", "actual:", result)
			}

			result = calcRealOdate(current, "+002", schedule)
			if result != "20201107" {
				t.Error("unexpected error, expecting odate:", "20201107", "actual:", result)
			}
		}
	}

func TestGetFirstDay(t *testing.T) {

	dd := date.Odate("20210401")
	getStartOfWeek(dd, 0)

}

func Test_prepareNextCycle_From_Start(t *testing.T) {

		interval := 2
		builder := taskdef.DummyTaskBuilder{}
		definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
			WithCyclic(taskdef.CyclicTaskData{
				IsCycle:      true,
				MaxRuns:      3,
				RunFrom:      taskdef.CycleFromStart,
				TimeInterval: interval,
			}).WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDaily,
		}).Build()

		if err != nil {
			t.Log(err)
		}

		task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

		tm := task.SetStartTime()
		expect := tm.Add(time.Duration(interval) * time.Minute)
		task.SetEndTime()

		task.prepareNextCycle()

		if task.CycleData().NextRun != types.FromTime(expect) {
			t.Error("unexpected value:", task.CycleData().NextRun, "expected:", types.FromTime(expect))
		}
	}

func Test_prepareNextCycle_From_End(t *testing.T) {

		interval := 2
		builder := taskdef.DummyTaskBuilder{}
		definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
			WithCyclic(taskdef.CyclicTaskData{
				IsCycle:      true,
				MaxRuns:      3,
				RunFrom:      taskdef.CycleFromEnd,
				TimeInterval: interval,
			}).WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDaily,
		}).Build()

		if err != nil {
			t.Log(err)
		}

		task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

		task.SetStartTime()
		tm := task.SetEndTime()
		expect := tm.Add(time.Duration(interval) * time.Minute)

		task.prepareNextCycle()

		if task.CycleData().NextRun != types.FromTime(expect) {
			t.Error("unexpected value:", task.CycleData().NextRun, "expected:", types.FromTime(expect))
		}
	}

func Test_prepareNextCycle_From_Schedule(t *testing.T) {

		interval := 2
		builder := taskdef.DummyTaskBuilder{}
		definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
			WithCyclic(taskdef.CyclicTaskData{
				IsCycle:      true,
				MaxRuns:      3,
				RunFrom:      taskdef.CycleFromSchedule,
				TimeInterval: interval,
			}).WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDaily,
		}).Build()

		if err != nil {
			t.Log(err)
		}

		task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

		task.SetStartTime()
		tm := task.SetEndTime()
		expect := tm.Add(time.Duration(interval) * time.Minute)

		task.prepareNextCycle()

		if task.CycleData().NextRun != types.FromTime(expect) {
			t.Error("unexpected value:", task.CycleData().NextRun, "expected:", types.FromTime(expect))
		}
	}

func Test_prepareNextCycle_MaxRuns(t *testing.T) {

		interval := 2
		maxRuns := 3
		builder := taskdef.DummyTaskBuilder{}
		definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
			WithCyclic(taskdef.CyclicTaskData{
				IsCycle:      true,
				MaxRuns:      maxRuns,
				RunFrom:      taskdef.CycleFromStart,
				TimeInterval: interval,
			}).WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDaily,
		}).Build()

		if err != nil {
			t.Log(err)
		}

		task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

		for i := 1; i <= maxRuns; i++ {

			task.SetStartTime()
			task.SetEndTime()
			result := task.prepareNextCycle()
			if !result && i < maxRuns {
				t.Error("unexpected value:", result, "expected:", true)
			}

			if result && i == maxRuns {
				t.Error("unexpected value:", result, "expected:", false)
			}

			task.SetExecutionID()
		}
	}
*/
func strTimeToInt(time string) (int, int) {
	val := strings.Split(time, ":")
	h, _ := strconv.Atoi(val[0])
	m, _ := strconv.Atoi(val[1])
	return h, m

}

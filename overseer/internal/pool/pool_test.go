package pool

/*
type mockResourceManager struct {
}

func (m *mockResourceManager) CheckTickets(in []types.CollectedTicketModel) []types.CollectedTicketModel {
	return []types.CollectedTicketModel{}
}
func (m *mockResourceManager) ProcessTicketAction([]types.TicketActionModel) bool {
	return false
}
func (m *mockResourceManager) ProcessAcquireFlag([]types.FlagModel) (bool, []string) {
	return false, []string{}
}
func (m *mockResourceManager) ProcessReleaseFlag([]string) (bool, []string) {
	return false, []string{}
}

type mockJournal struct {
	timeout   int
	collected chan events.RouteJournalMsg
}

func (j *mockJournal) Push(msg events.RouteJournalMsg) {
	j.collected <- msg

}

func (j *mockJournal) Collect(expected int, after time.Time) <-chan []events.RouteJournalMsg {

	ch := make(chan []events.RouteJournalMsg)

	go func(timeout, expected int, done chan<- []events.RouteJournalMsg, col chan events.RouteJournalMsg) {

		collected := 0
		deadline := time.After(time.Duration(timeout) * time.Second)
		result := []events.RouteJournalMsg{}
		for {
			select {
			case <-deadline:
				{
					close(done)
					return
				}
			case d := <-col:
				{
					if d.Time.Before(after) {
						continue
					}
					result = append(result, d)
					collected++
					if collected == expected {
						done <- result
						close(done)
						return
					}

				}
			}
		}

	}(j.timeout, expected, ch, j.collected)

	return ch

}

type mockWorkerManager struct {
}

type mockJournalWriter struct {
}

func (m *mockJournalWriter) PushJournalMessage(ID unique.TaskOrderID, execID string, t time.Time, msg string) {

}

func (m *mockWorkerManager) Push(ctx context.Context, t types.TaskDescription, vars types.EnvironmentVariableList) (types.WorkerTaskStatus, error) {
	return types.WorkerTaskStatusRecieved, nil
}
func (m *mockWorkerManager) Status(ctx context.Context, t types.WorkDescription) types.TaskExecutionStatus {
	return types.TaskExecutionStatus{}
}

const (
	testCollectionName = "tasks"
	testStoreTaskName  = "storetasks"
	testSequenceName   = "sequence"
)

var storeConfig config.StoreConfiguration = config.StoreConfiguration{ID: "teststore", ConnectionString: "local;/../../../data/tests"}

var taskPoolConfig config.ActivePoolConfiguration = config.ActivePoolConfiguration{
	ForceNewDayProc: true, MaxOkReturnCode: 4,
	NewDayProc: "00:30",
	SyncTime:   5,
}

type mockSequence struct {
	val int
}

func (m *mockSequence) Next() unique.TaskOrderID {

	m.val++
	return unique.TaskOrderID(fmt.Sprintf("%05d", m.val))
}

var seq = &mockSequence{val: 1}

var provider *datastore.Provider

var definitionManagerT taskdef.TaskDefinitionManager
var taskPoolT *ActiveTaskPool
var activeTaskManagerT *ActiveTaskPoolManager

var mockJournalT = &mockJournal{timeout: 3, collected: make(chan events.RouteJournalMsg, 10)}
var defManagerDircetory string
var isInitialized bool = false

func setupEnv() {

	f, _ := os.Create(fmt.Sprintf("../../../data/tests/%s.json", testCollectionName))
	f.Write([]byte("{}"))
	f.Close()

	f1, _ := os.Create(fmt.Sprintf("../../../data/tests/%s.json", testStoreTaskName))
	f1.Write([]byte("{}"))
	f1.Close()

	f2, _ := os.Create(fmt.Sprintf("../../../data/tests/%s.json", testSequenceName))
	f2.Write([]byte(`{}`))
	f2.Close()

	provider, _ = datastore.NewDataProvider(storeConfig)

	defManagerDircetory, _ = filepath.Abs("../../../def_test/")
	definitionManagerT, _ = taskdef.NewManager(defManagerDircetory)

	initTaskPool(provider)
	activeTaskManagerT, _ = NewActiveTaskPoolManager(definitionManagerT, taskPoolT, provider)
	//activeTaskManagerT.log = log
	activeTaskManagerT.sequence = seq

	isInitialized = true

}
func initTaskPool(prov *datastore.Provider) {

	taskPoolT, _ = NewTaskPool(taskPoolConfig, prov, true, definitionManagerT, &mockWorkerManager{}, &mockResourceManager{}, &mockJournalWriter{})
}
*/

package taskdef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"overseer/common/logger"
	"overseer/overseer/internal/unique"
	"overseer/overseer/taskdata"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	ErrLockIDNotExists  error = errors.New("given lockID does not exists")
	ErrTaskNameEmpty    error = errors.New("task name cannot be empty")
	ErrTaskRevEmpty     error = errors.New("task rev cannot be empty")
	ErrTaskRename       error = errors.New("unable to rename, definition already exists")
	ErrTaskRevDiff      error = errors.New("task rev is different than original")
	ErrTaskRevInvalid   error = errors.New("invalid task rev")
	ErrGroupNameEmpty   error = errors.New("group name cannot be empty")
	ErrGroupNotExists   error = errors.New("group with given name not exists")
	ErrGroupDirNotEmpty error = errors.New("group directory not empty, remove definitions firs")
	ErrGroupDirInvalid  error = errors.New("failed to read groups from directory,invalid path")
	ErrReferenceEmpty   error = errors.New("empty reference id")
)

var poolDirectoryName = ".pool"

type taskManager struct {
	dirPath       string
	activeTaskDir string
	lock          sync.Mutex
	log           logger.AppLogger
}

//TaskDefinitionManager - main component responsible for a task CRUD
type TaskDefinitionManager interface {
	GetTasks(tasks ...taskdata.GroupNameData) []TaskDefinition
	GetTask(task taskdata.GroupNameData) (TaskDefinition, error)
	GetGroups() ([]string, error)
	GetTasksFromGroup(groups []string) ([]taskdata.GroupNameData, error)
	GetTaskModelList(group taskdata.GroupData) ([]taskdata.TaskNameModel, error)
	Create(task TaskDefinition) error
	CreateGroup(name string) error
	Update(task TaskDefinition) error
	Delete(task taskdata.GroupNameData) error
	DeleteGroup(name string) error
	GetActiveDefinition(refID string) (TaskDefinition, error)
	WriteActiveDefinition(def TaskDefinition, id unique.MsgID) error
	RemoveActiveDefinition(id string) error
}

//NewManager - returns new instance of a TaskDefinitionManager
func NewManager(path string, log logger.AppLogger) (TaskDefinitionManager, error) {

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	poolDir := filepath.Join(path, poolDirectoryName)

	if _, err := os.Stat(poolDir); err != nil {
		return nil, err
	}

	var t = &taskManager{}

	t.dirPath = path
	t.activeTaskDir = poolDirectoryName
	t.log = log

	return t, nil
}
func (m *taskManager) GetTask(task taskdata.GroupNameData) (TaskDefinition, error) {

	var err error
	var result TaskDefinition

	defPath := filepath.Join(m.dirPath, task.Group, fmt.Sprintf("%v.json", task.Name))
	if result, err = FromDefinitionFile(defPath); err != nil {
		return nil, err
	}

	return result, nil
}
func (m *taskManager) GetTasks(tasks ...taskdata.GroupNameData) []TaskDefinition {

	result := make([]TaskDefinition, 0)
	for _, n := range tasks {

		if t, err := m.GetTask(n); err == nil {
			result = append(result, t)
		} else {
			m.log.Error("GetTasks:", err)
		}

	}
	return result
}

func (m *taskManager) GetTasksFromGroup(groups []string) ([]taskdata.GroupNameData, error) {

	result := make([]taskdata.GroupNameData, 0)

	for _, grp := range groups {

		pth := filepath.Join(m.dirPath, grp)
		info, err := ioutil.ReadDir(pth)
		if err != nil {
			return nil, errors.New("can't find group with given name")
		}

		for _, nfo := range info {

			if !nfo.IsDir() {

				nameExt := strings.Split(nfo.Name(), ".")
				result = append(result, taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: grp}, Name: nameExt[0]})
			}
		}
	}

	return result, nil
}
func (m *taskManager) GetTaskModelList(group taskdata.GroupData) ([]taskdata.TaskNameModel, error) {

	var info []os.FileInfo
	var err error
	var definition TaskDefinition
	result := []taskdata.TaskNameModel{}

	pth := filepath.Join(m.dirPath, group.Group)
	if info, err = ioutil.ReadDir(pth); err != nil {
		return nil, err
	}

	for _, nfo := range info {

		if !nfo.IsDir() {
			fpath := filepath.Join(pth, nfo.Name())
			if definition, err = FromDefinitionFile(fpath); err != nil {
				continue
			}
			n, g, d := definition.GetInfo()
			result = append(result, taskdata.TaskNameModel{Group: g, Name: n, Description: d})
		}
	}

	return result, nil
}

func (m *taskManager) GetGroups() ([]string, error) {

	groups := make([]string, 0)

	info, err := ioutil.ReadDir(m.dirPath)
	if err != nil {
		m.log.Error("Get groups:", err)
		return []string{}, ErrGroupDirInvalid
	}

	for _, in := range info {
		if in.IsDir() && !strings.HasPrefix(in.Name(), ".") {
			groups = append(groups, in.Name())
		}
	}

	return groups, nil
}
func (m *taskManager) Create(task TaskDefinition) error {
	defer m.lock.Unlock()
	m.lock.Lock()

	nm, grp, _ := task.GetInfo()
	path := filepath.Join(m.dirPath, grp, fmt.Sprintf("%v.json", nm))

	_, err := os.Stat(path)
	if err == nil {
		return errors.New("unable to create, definition already exists")
	}

	if task.Rev() == "" {
		task.SetRevision(unique.NewID())
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0640)

	return err
}

func (m *taskManager) CreateGroup(name string) error {

	if name == "" {
		return ErrGroupNameEmpty
	}
	err := os.Mkdir(filepath.Join(m.dirPath, name), os.ModeDir|0640)
	if err != nil {
		return err
	}
	return nil
}

func (m *taskManager) Update(task TaskDefinition) error {
	defer m.lock.Unlock()
	m.lock.Lock()

	var err error
	name, grp, _ := task.GetInfo()

	if name == "" {
		return ErrTaskNameEmpty
	}

	n, g, _, err := getNameGroupIdFromDefinition(task)
	if err != nil {
		return err
	}

	opath := filepath.Join(m.dirPath, g, fmt.Sprintf("%v.json", n))

	oldDef, err := FromDefinitionFile(opath)
	if err != nil {
		return err
	}
	//ummy_update_03@test@16be1f583874c7cf61b6be1f
	//ummy_update_03@test@16be1f583874c7cf61b6be1f

	if oldDef.Rev() != task.Rev() {
		return ErrTaskRevDiff
	}

	path := filepath.Join(m.dirPath, grp, fmt.Sprintf("%v.json", name))

	//only if task is moved to different group or renamed, ensure if task with new name or new group does not exists
	if name != n || grp != g {
		path := filepath.Join(m.dirPath, grp, fmt.Sprintf("%v.json", name))
		_, err = os.Stat(path)

		if err == nil {
			return ErrTaskRename
		}
	}

	var result string

	if result, err = SerializeDefinition(task); err != nil {
		return err
	}

	if path != opath {
		os.Remove(opath)
	}

	return os.WriteFile(path, []byte(result), 0640)

}
func (m *taskManager) Delete(task taskdata.GroupNameData) error {

	defer m.lock.Unlock()
	m.lock.Lock()

	path := filepath.Join(m.dirPath, task.Group, fmt.Sprintf("%v.json", task.Name))
	return os.Remove(path)
}

func (m *taskManager) DeleteGroup(name string) error {
	defer m.lock.Unlock()
	m.lock.Lock()

	if name == "" {
		return ErrGroupNameEmpty
	}

	path := filepath.Join(m.dirPath, name)
	finfo, err := os.Open(path)
	if err != nil {
		return ErrGroupNotExists
	}

	var stat os.FileInfo
	var errs error
	stat, errs = finfo.Stat()
	if errs != nil {
		return errs
	}
	if !stat.IsDir() {
		return errors.New("given name is not a directory")
	}

	_, err = finfo.Readdirnames(1)
	if err == nil {
		return ErrGroupDirNotEmpty
	}
	if err == io.EOF {
		finfo.Close()
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *taskManager) WriteActiveDefinition(task TaskDefinition, id unique.MsgID) error {

	var data []byte
	var err error
	path := filepath.Join(m.dirPath, m.activeTaskDir, fmt.Sprintf("%s.json", id.Hex()))

	if data, err = json.Marshal(task); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
func (m *taskManager) RemoveActiveDefinition(id string) error {

	path := filepath.Join(m.dirPath, m.activeTaskDir, fmt.Sprintf("%s.json", id))
	return os.Remove(path)
}
func (m *taskManager) GetActiveDefinition(id string) (TaskDefinition, error) {

	var err error
	var result TaskDefinition

	if id == "" {
		return nil, ErrReferenceEmpty
	}

	path := filepath.Join(m.dirPath, m.activeTaskDir, fmt.Sprintf("%s.json", id))
	if result, err = FromPoolDirectory(path); err != nil {
		return nil, err
	}

	return result, nil
}

//getNameGroupFrovRev - gets name of a task, group that belongs and unique idetifier,
//returns error if revision is not  valid format which is: task_name@task_group@unique_id
func getNameGroupIdFromDefinition(task TaskDefinition) (string, string, string, error) {

	rev := task.Rev()

	if rev == "" {
		return "", "", "", ErrTaskRevEmpty
	}

	m, err := regexp.Match(`^[A-Za-z][\w\-\.]*@[A-Za-z][\w\-\.]*@[\w]+$`, []byte(rev))
	if err != nil || !m {
		return "", "", "", ErrTaskRevInvalid
	}

	val := strings.Split(rev, "@")

	return val[0], val[1], val[2], nil
}

package taskdef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"overseer/common/logger"
	"overseer/overseer/taskdata"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	ErrLockIDNotExists error  = errors.New("given lockID does not exists")
	ErrTaskNameEmpty   error  = errors.New("task name cannot be empty")
	ErrGroupNameEmpty  error  = errors.New("group name cannot be empty")
	lockIDSeq          uint32 = 1
)

type lockData struct {
	group string
	name  string
	file  *os.File
}

type taskManager struct {
	dirPath string
	lockTab map[uint32]lockData
	lock    sync.Mutex
	log     logger.AppLogger
}

//TaskDefinitionManager - main component responsible for a task CRUD
type TaskDefinitionManager interface {
	GetTasks(tasks ...taskdata.GroupNameData) []TaskDefinition
	GetTask(task taskdata.GroupNameData) (TaskDefinition, error)
	GetGroups() []string
	GetTasksFromGroup(groups []string) ([]taskdata.GroupNameData, error)
	GetTaskModelList(filter string) ([]taskdata.TaskNameModel, error)
	Lock(task taskdata.GroupNameData) (uint32, error)
	Unlock(lockID uint32) error
	Create(task TaskDefinition) error
	CreateGroup(name string) error
	Update(lockID uint32, task TaskDefinition) error
	Delete(lockID uint32, task taskdata.GroupNameData) error
	DeleteGroup(name string) error
}

//NewManager - returns new instance of a TaskDefinitionManager
func NewManager(path string) (TaskDefinitionManager, error) {

	var t = new(taskManager)
	t.dirPath = path
	t.lockTab = make(map[uint32]lockData, 0)
	t.log = logger.Get()

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
			return nil, errors.New("Can't find group with given name")
		}

		for _, nfo := range info {

			if nfo.IsDir() == false {

				nameExt := strings.Split(nfo.Name(), ".")
				result = append(result, taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: grp}, Name: nameExt[0]})
			}
		}
	}

	return result, nil
}
func (m *taskManager) GetTaskModelList(filter string) ([]taskdata.TaskNameModel, error) {

	var info []os.FileInfo
	var err error
	var definition TaskDefinition
	result := []taskdata.TaskNameModel{}

	pth := filepath.Join(m.dirPath, filter)
	if info, err = ioutil.ReadDir(pth); err != nil {
		return nil, err
	}

	for _, nfo := range info {

		if nfo.IsDir() == false {
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

func (m *taskManager) GetGroups() []string {

	groups := make([]string, 0)

	info, err := ioutil.ReadDir(m.dirPath)
	if err != nil {
		m.log.Error("GetTasks:", err)
	}

	for _, in := range info {
		if in.IsDir() {
			groups = append(groups, in.Name())
		}
	}

	return groups
}
func (m *taskManager) Lock(task taskdata.GroupNameData) (uint32, error) {
	defer m.lock.Unlock()
	m.lock.Lock()

	if task.Name == "" {
		return 0, ErrTaskNameEmpty
	}

	for _, n := range m.lockTab {

		if n.group == task.Group && n.name == task.Name {
			return 0, errors.New("Unable to acquire lock")
		}
	}
	f, err := os.Open(filepath.Join(m.dirPath, task.Group, fmt.Sprintf("%v.json", task.Name)))
	if err != nil {
		return 0, err
	}

	lockID := getNext()
	m.lockTab[lockID] = lockData{name: task.Name, group: task.Group, file: f}

	return lockID, nil
}
func (m *taskManager) Unlock(lockID uint32) error {
	defer m.lock.Unlock()
	m.lock.Lock()
	d, x := m.lockTab[lockID]
	if x == false {
		return ErrLockIDNotExists
	}
	d.file.Close()

	delete(m.lockTab, lockID)
	return nil

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
func (m *taskManager) Update(lockID uint32, task TaskDefinition) error {
	defer m.lock.Unlock()
	m.lock.Lock()

	d, x := m.lockTab[lockID]
	if x == false {
		return ErrLockIDNotExists
	}

	name, grp, _ := task.GetInfo()

	if name == "" {
		return ErrTaskNameEmpty
	}
	//ensure also if task with new name or new group does not exists
	if d.name != name || d.group != grp {
		path := filepath.Join(m.dirPath, grp, fmt.Sprintf("%v.json", name))
		oldpath := filepath.Join(m.dirPath, d.group, fmt.Sprintf("%v.json", d.name))
		_, err := os.Stat(path)

		if err == nil {
			return errors.New("unable to rename, definition already exists")
		}

		d.file.Close()
		err = os.Rename(oldpath, path)
		nfile, err := os.OpenFile(path, os.O_TRUNC|os.O_RDWR, 0640)

		m.lockTab[lockID] = lockData{file: nfile, name: name, group: grp}
		d, _ = m.lockTab[lockID]

	}

	var result string
	var err error
	result, err = SerializeDefinition(task)
	if err != nil {
		return err
	}

	_, err = d.file.Seek(0, 0)
	if err != nil {
		return err
	}

	if err = d.file.Truncate(0); err != nil {
		return err
	}
	if _, err = d.file.Write([]byte(result)); err != nil {
		return err
	}

	return nil
}
func (m *taskManager) Delete(lockID uint32, task taskdata.GroupNameData) error {

	defer m.lock.Unlock()
	m.lock.Lock()

	d, x := m.lockTab[lockID]
	if x == false {
		return ErrLockIDNotExists
	}
	if task.Name != d.name || task.Group != d.group {
		return errors.New("group and name does not match with lockID")
	}

	path := filepath.Join(m.dirPath, task.Group, fmt.Sprintf("%v.json", task.Name))
	d.file.Close()
	os.Remove(path)
	delete(m.lockTab, lockID)
	return nil
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
		return errors.New("can't find directory")
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
		return errors.New("directory is not empty")
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

func getNext() uint32 {

	return atomic.AddUint32(&lockIDSeq, 1)
}

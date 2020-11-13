package resources

import (
	"encoding/json"
	"fmt"
	"goscheduler/overseer/internal/date"
	"os"
	"path/filepath"
)

type (

	//FlagResourcePolicy - type of flag
	FlagResourcePolicy int8

	//TicketResource - Condition resources
	TicketResource struct {
		Name  string     `validate:"required,max=32"`
		Odate date.Odate `validate:"required,odate"`
	}
	//FlagResource - Semaphore like resources
	FlagResource struct {
		Name   string `validate:"required,max=32"`
		Policy FlagResourcePolicy
		Count  int
	}
)

const (
	//FlagPolicyShared  - task can run together with other tasks that share this resources
	FlagPolicyShared FlagResourcePolicy = 0
	//FlagPolicyExclusive - only one task can run with exclusive policy
	FlagPolicyExclusive FlagResourcePolicy = 1
)

type resourcePool struct {
	fpath   string
	tpath   string
	Tickets []TicketResource `json:"tickets"`
	Flags   []FlagResource   `json:"flags"`
}

//newResourcePool - Creates a new pool that holds resources.
func newResourcePool(directory string) (*resourcePool, error) {

	absDir, _ := filepath.Abs(directory)
	pool := &resourcePool{Tickets: make([]TicketResource, 0), Flags: make([]FlagResource, 0)}
	pool.tpath = filepath.Join(absDir, "tickets.json")
	pool.fpath = filepath.Join(absDir, "flags.json")

	data, err := getFileContent(pool.tpath)
	if err != nil {
		return pool, err
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &pool.Tickets)
		if err != nil {
			fmt.Println(err)
			return pool, err
		}
	}

	data, err = getFileContent(pool.fpath)
	if err != nil {
		fmt.Println(err)
		return pool, err
	}
	if len(data) != 0 {
		err = json.Unmarshal(data, &pool.Flags)
		if err != nil {
			fmt.Println(err)
			return pool, err
		}
	}

	return pool, nil

}

func getFileContent(path string) ([]byte, error) {
	//:TODO
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	stat, _ := f.Stat()
	sz := stat.Size()
	data := make([]byte, sz)
	f.Seek(0, 0)
	_, err = f.Read(data)
	return data, err

}

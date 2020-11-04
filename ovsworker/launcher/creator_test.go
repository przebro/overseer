package launcher

import (
	"goscheduler/common/types"
	"goscheduler/ovsworker/msgheader"
	"testing"
)

func TestCreateFactory(t *testing.T) {

	l := &FragmentLauncher{}
	creator := FragmentFactory(l)
	if creator == nil {
		t.Error("create FragmentFactory error")
	}

}

func TestCreateFragment(t *testing.T) {
	l := &FragmentLauncher{store: NewStore()}
	creator := FragmentFactory(l)
	if creator == nil {
		t.Error("create FragmentFactory error")
	}
	header := msgheader.TaskHeader{Type: types.TypeDummy, TaskID: "12345", Variables: make(map[string]string, 0)}
	data := make([]byte, 0)
	err := creator.CreateFragment(header, data)
	if err != nil {
		t.Error(err)
	}
}

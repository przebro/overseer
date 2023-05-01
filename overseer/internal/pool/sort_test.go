package pool

import (
	"testing"

	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/events"
)

func TestSortSwap(t *testing.T) {
	sorter := taskInfoSorter{list: make([]events.TaskInfoResultMsg, 0)}
	sorter.list = append(sorter.list, events.TaskInfoResultMsg{Name: "TEST01"},
		events.TaskInfoResultMsg{Name: "TEST02"},
		events.TaskInfoResultMsg{Name: "TEST03"},
	)

	sorter.Swap(0, 2)
	if sorter.list[0].Name != "TEST03" && sorter.list[2].Name != "TEST01" {
		t.Error("unexpected result")

	}

}

func TestSortLess(t *testing.T) {
	sorter := taskInfoSorter{list: make([]events.TaskInfoResultMsg, 0)}
	sorter.list = append(sorter.list, events.TaskInfoResultMsg{Name: "TEST01", TaskID: unique.TaskOrderID("00001")},
		events.TaskInfoResultMsg{Name: "TEST02", TaskID: unique.TaskOrderID("00010")},
		events.TaskInfoResultMsg{Name: "TEST03", TaskID: unique.TaskOrderID("00002")},
	)

	result := sorter.Less(0, 1)
	if !result {
		t.Error("unexpected result")

	}

	result = sorter.Less(1, 0)
	if result {
		t.Error("unexpected result")

	}

}

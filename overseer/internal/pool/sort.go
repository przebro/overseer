package pool

import (
	"github.com/przebro/overseer/overseer/internal/events"
)

type taskInfoSorter struct{ list []events.TaskInfoResultMsg }

func (s taskInfoSorter) Len() int      { return len(s.list) }
func (s taskInfoSorter) Swap(i, j int) { s.list[i], s.list[j] = s.list[j], s.list[i] }

func (s taskInfoSorter) Less(i, j int) bool {
	return s.list[i].TaskID < s.list[j].TaskID
}

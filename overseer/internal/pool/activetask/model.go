package activetask

import (
	"time"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/internal/pool/models"
)

type ActiveTaskModel struct {
	OrderID    string               `json:"_id" bson:"_id"`
	Name       string               `json:"name" bson:"name"`
	Group      string               `json:"group" bson:"group"`
	Reference  string               `json:"_ref" bson:"_ref"`
	State      models.TaskState     `json:"state" bson:"state"`
	Holded     bool                 `json:"holded" bson:"holded"`
	Confirmed  bool                 `json:"confirm" bson:"confirm"`
	OrderDate  date.Odate           `json:"odate" bson:"odate"`
	Tickets    []taskInTicketModel  `json:"tickets" bson:"tickets"`
	RunNumber  int32                `json:"rn" bson:"rn"`
	Executions []taskExecutionModel `json:"exec" bson:"exec"`
	Cycle      taskCycleModel       `json:"cycle" bson:"cycle"`
	Waiting    string               `json:"waiting" bson:"waiting"`
}

type taskInTicketModel struct {
	Name      string `json:"name" bson:"name"`
	Odate     string `json:"odate" bson:"odate"`
	Fulfilled bool   `json:"ff" bson:"ff"`
}

type taskExecutionModel struct {
	ID        string           `json:"_id" bson:"_id"`
	Worker    string           `json:"worker,omitempty" bson:"worker,omitempty"`
	StartTime time.Time        `json:"start,omitempty" bson:"start,omitempty"`
	EndTime   time.Time        `json:"end,omitempty" bson:"end,omitempty"`
	State     models.TaskState `json:"state" bson:"state"`
}

type TaskPoolModel struct {
	ID   string            `json:"_id" bson:"_id"`
	Data []ActiveTaskModel `json:"data" bson:"data"`
}

type taskCycleModel struct {
	IsCyclic bool   `json:"is" bson:"is"`
	NextRun  string `json:"tm" bson:"tm"`
	MaxRun   int    `json:"max" bson:"max"`
	RunFrom  string `json:"rf" bson:"rf"`
	Interval int    `json:"in" bson:"in"`
}

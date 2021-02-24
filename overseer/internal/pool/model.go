package pool

import (
	"encoding/json"
	"overseer/common/types/date"
	"time"
)

type activeTaskModel struct {
	OrderID    string               `json:"_id" bson:"_id"`
	Definition json.RawMessage      `json:"def" bson:"def"`
	State      TaskState            `json:"state" bson:"state"`
	Holded     bool                 `json:"holded" bson:"holded"`
	Confirmed  bool                 `json:"confirm" bson:"confirm"`
	OrderDate  date.Odate           `json:"odate" bson:"odate"`
	Tickets    []taskInTicketModel  `json:"tickets" bson:"tickets"`
	RunNumber  int32                `json:"rn" bson:"rn"`
	Executions []taskExecutionModel `json:"exec" bson:"exec"`
	Waiting    string               `json:"waiting" bson:"waiting"`
}

type taskInTicketModel struct {
	Name      string `json:"name" bson:"name"`
	Odate     string `json:"odate" bson:"odate"`
	Fulfilled bool   `json:"ff" bson:"ff"`
}

type taskExecutionModel struct {
	ID        string    `json:"_id" bson:"_id"`
	Worker    string    `json:"worker,omitempty" bson:"worker,omitempty"`
	StartTime time.Time `json:"start,omitempty" bson:"start,omitempty"`
	EndTime   time.Time `json:"end,omitempty" bson:"end,omitempty"`
	State     TaskState `json:"state" bson:"state"`
}

type taskPoolModel struct {
	ID   string            `json:"_id" bson:"_id"`
	Data []activeTaskModel `json:"data" bson:"data"`
}

package pool

import (
	"encoding/json"
	"overseer/common/types/date"
	"time"
)

type activeTaskModel struct {
	OrderID    string              `json:"_id" bson:"_id"`
	Definition json.RawMessage     `json:"definition" bson:"definition"`
	State      TaskState           `json:"state" bson:"state"`
	Holded     bool                `json:"holded" bson:"holded"`
	Confirmed  bool                `json:"confirm" bson:"confirm"`
	OrderDate  date.Odate          `json:"odate" bson:"odate"`
	Tickets    []taskInTicketModel `json:"tickets" bson:"tickets"`
	RunNumber  int32               `json:"rn" bson:"rn"`
	Worker     string              `json:"worker" bson:"worker"`
	Waiting    string              `json:"waiting" bson:"waiting"`
	StartTime  time.Time           `json:"start" bson:"start"`
	EndTime    time.Time           `json:"end" bson:"end"`
	Output     []string            `json:"out" bson:"out"`
}

type taskInTicketModel struct {
	Name      string `json:"name" bson:"name"`
	Odate     string `json:"odate" bson:"odate"`
	Fulfilled bool   `json:"ff" bson:"ff"`
}

type taskPoolModel struct {
	ID   string            `json:"_id" bson:"_id"`
	Data []activeTaskModel `json:"data" bson:"data"`
}

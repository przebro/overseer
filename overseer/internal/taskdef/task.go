package taskdef

import (
	"encoding/json"

	"github.com/przebro/overseer/common/types"
)

type baseTaskDefinition_DEPRECATED struct {
	Revision      string                        `json:"rev"`
	TaskType      types.TaskType                `json:"type" validate:"oneof=dummy os aws filewatch"`
	Name          string                        `json:"name" validate:"required,max=32,resname"`
	Group         string                        `json:"group" validate:"required,max=20,resname"`
	Description   string                        `json:"description" validate:"lte=200"`
	ConfirmFlag   bool                          `json:"confirm"`
	DataRetention int                           `json:"retention" validate:"min=0,max=14"`
	Schedule      SchedulingData                `json:"schedule" validate:"omitempty"`
	Cyclics       CyclicTaskData                `json:"cyclic" validate:"omitempty"`
	InRelation    InTicketRelation              `json:"relation" validate:"required_with=InTickets,omitempty,oneof=AND OR EXPR"`
	Expression    string                        `json:"expr" validate:"omitempty"`
	InTickets     []InTicketData                `json:"inticket" validate:"omitempty,dive"`
	FlagsTab      []FlagData                    `json:"flags"  validate:"omitempty,dive"`
	OutTickets    []OutTicketData               `json:"outticket"  validate:"omitempty,dive"`
	TaskVariables types.EnvironmentVariableList `json:"variables"  validate:"omitempty,dive"`
	Data          json.RawMessage               `json:"spec,omitempty"`
	Worker        string                        `json:"worker,omitempty"`
}

package resources

import (
	"overseer/overseer/internal/date"
)

type (

	//FlagResourcePolicy - type of flag
	FlagResourcePolicy int8

	//TicketResource - Condition resources
	TicketResource struct {
		Name  string     `json:"name" bson:"name" validate:"required,max=32"`
		Odate date.Odate `json:"odate" bson:"odate"`
	}
	//FlagResource - Semaphore like resources
	FlagResource struct {
		Name   string             `json:"name" bson:"name validate:"required,max=32"`
		Policy FlagResourcePolicy `json:"policy" bson:"policy"`
		Count  int                `json:"count" bson:"count"`
	}

	TicketsResourceModel struct {
		ID      string           `json:"_id" bson:"_id`
		REV     string           `json:"_rev" bson:"_rev"`
		Tickets []TicketResource `json:"tickets" bson:"tickets"`
	}

	FlagsResourceModel struct {
		ID    string         `json:"_id" bson:"_id`
		REV   string         `json:"_rev" bson:"_rev"`
		Flags []FlagResource `json:"flags" bson:"flags"`
	}
)

const (
	//FlagPolicyShared  - task can run together with other tasks that share this resources
	FlagPolicyShared FlagResourcePolicy = 0
	//FlagPolicyExclusive - only one task can run with exclusive policy
	FlagPolicyExclusive FlagResourcePolicy = 1
)

type ticketSorter struct{ list []TicketResource }

func (s ticketSorter) Len() int      { return len(s.list) }
func (s ticketSorter) Swap(i, j int) { s.list[i], s.list[j] = s.list[j], s.list[i] }

func (s ticketSorter) Less(i, j int) bool {
	return s.list[i].Name < s.list[j].Name
}

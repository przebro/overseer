package resources

type (
	ResourceType string

	ResourceModel struct {
		ID    string       `json:"_id" bson:"_id"`
		Type  ResourceType `json:"type" bson:"type"`
		Value int64        `json:"value"`
	}

	FlagFilterOptions struct {
		FlagPolicy []uint8
	}
	TicketFilterOptions struct {
		Odate         string
		OrderDateFrom string
		OrderDateTo   string
	}

	ResourceFilter struct {
		Name string
		Type []ResourceType
		*FlagFilterOptions
		*TicketFilterOptions
	}
)

const (
	//FlagPolicyShared  - task can run together with other tasks that share this resources
	FlagPolicyShared uint8 = 0
	//FlagPolicyExclusive - only one task can run with exclusive policy
	FlagPolicyExclusive uint8 = 1

	ResourceTypeTicket ResourceType = "ticket"
	ResourceTypeFlag   ResourceType = "flag"
)

type ticketSorter struct{ list []ResourceModel }

func (s ticketSorter) Len() int      { return len(s.list) }
func (s ticketSorter) Swap(i, j int) { s.list[i], s.list[j] = s.list[j], s.list[i] }

func (s ticketSorter) Less(i, j int) bool {
	return s.list[i].ID < s.list[j].ID
}

package types

import "github.com/przebro/overseer/common/types/date"

type CollectedTicketModel struct {
	Name   string
	Odate  date.Odate
	Exists bool
}

type OutAction string

type TicketActionModel struct {
	Name   string
	Odate  date.Odate
	Action OutAction
}

type FlagModel struct {
	Name   string
	Policy uint8
}

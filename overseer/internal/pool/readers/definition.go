package readers

import (
	"time"

	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/taskdef"
)

// ActiveDefinitionReader - reads definitions
type ActiveDefinitionReader interface {
	GetActiveDefinition(refID string) (*taskdef.TaskDefinition, error)
}

// ActiveDefinitionWriter - writes definitions
type ActiveDefinitionWriter interface {
	WriteActiveDefinition(def *taskdef.TaskDefinition, id unique.MsgID) error
}

// ActiveDefinitionRemover - removes definitions
type ActiveDefinitionRemover interface {
	RemoveActiveDefinition(id string) error
}

// ActiveDefinitionReadWriter - groups definition reader and writer
type ActiveDefinitionReadWriter interface {
	ActiveDefinitionReader
	ActiveDefinitionWriter
}

// ActiveDefinitionReadWriter - groups definition reader, writer and remover
type ActiveDefinitionReadWriterRemover interface {
	ActiveDefinitionReader
	ActiveDefinitionWriter
	ActiveDefinitionRemover
}

type JournalWriter interface {
	PushJournalMessage(ID unique.TaskOrderID, execID string, t time.Time, msg string)
}

package work

import (
	"overseer/common/logger"
	"overseer/overseer/internal/taskdef"
	"overseer/proto/actions"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

//ActionConverter -  Converts from an internal task to protobuf type Any.
type ActionConverter interface {
	Convert(data interface{}, variables []taskdef.VariableData, log logger.AppLogger) *any.Any
}

//DummyActionConverter - Chains link for an Dummy action type.
type DummyActionConverter struct {
	next ActionConverter
}

//OsTaskActionConverter - Chains link for an OS action type.
type OsTaskActionConverter struct {
	next ActionConverter
}

//NewConverterChain - Creates a new converter chain.
func NewConverterChain() ActionConverter {

	c := &DummyActionConverter{
		next: &OsTaskActionConverter{
			next: nil,
		},
	}

	return c
}

//Convert - Converts from an internal task to protobuf type Any.
func (c *DummyActionConverter) Convert(data interface{}, variables []taskdef.VariableData, log logger.AppLogger) *any.Any {

	result, isOk := data.(string)
	if isOk {
		cmd := &actions.DummyTaskAction{Data: result}
		if act, err := proto.Marshal(cmd); err == nil {
			msg := &any.Any{TypeUrl: string(cmd.ProtoReflect().Descriptor().FullName()), Value: act}
			return msg
		}
		log.Error("DummyActionConverter:unable to convert data")
		return nil

	}
	if c.next == nil {
		return nil
	}
	return c.next.Convert(data, variables, log)
}

//Convert - Converts from an internal task to protobuf type Any.
func (c *OsTaskActionConverter) Convert(data interface{}, variables []taskdef.VariableData, log logger.AppLogger) *any.Any {

	result, isOk := data.(taskdef.OsTaskData)
	if isOk {

		reg := regexp.MustCompile(`\%\%[A-Z0-9_]+`)
		cmdLine := result.CommandLine
		for _, n := range variables {
			if reg.MatchString(n.Name) {
				cmdLine = strings.Replace(cmdLine, n.Name, n.Value, -1)
			}
		}

		cmd := &actions.OsTaskAction{CommandLine: cmdLine, Runas: result.RunAs, Type: string(result.ActionType)}
		if act, err := proto.Marshal(cmd); err == nil {

			msg := &any.Any{TypeUrl: string(cmd.ProtoReflect().Descriptor().FullName()), Value: act}
			return msg
		}
		log.Error("OsActionConverter:unable to convert data")
		return nil
	}

	if c.next == nil {
		return nil
	}
	return c.next.Convert(data, variables, log)
}

package services

import (
	"context"
	"overseer/common/logger"
	"overseer/proto/services"

	"github.com/golang/protobuf/ptypes/empty"
)

type ovsStatusService struct {
	log logger.AppLogger
}

//NewStatusService - Creates a new status service
func NewStatusService() services.StatusServiceServer {
	return &ovsStatusService{log: logger.Get()}

}

//OverseerStatus - tests connection
func (srv *ovsStatusService) OverseerStatus(ctx context.Context, msg *empty.Empty) (*services.ActionResultMsg, error) {

	result := &services.ActionResultMsg{
		Success: true, Message: "ovs status service ready",
	}

	return result, nil
}

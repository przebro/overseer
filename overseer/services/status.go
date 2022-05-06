package services

import (
	"context"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/proto/services"

	empty "google.golang.org/protobuf/types/known/emptypb"
)

type ovsStatusService struct {
	log logger.AppLogger
	services.UnimplementedStatusServiceServer
}

//NewStatusService - Creates a new status service
func NewStatusService(log logger.AppLogger) services.StatusServiceServer {
	return &ovsStatusService{log: log}

}

//OverseerStatus - tests connection
func (srv *ovsStatusService) OverseerStatus(ctx context.Context, msg *empty.Empty) (*services.ActionResultMsg, error) {

	result := &services.ActionResultMsg{
		Success: true, Message: "ovs status service ready",
	}

	return result, nil
}

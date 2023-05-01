package services

import (
	"context"

	"github.com/przebro/overseer/proto/services"

	empty "google.golang.org/protobuf/types/known/emptypb"
)

type ovsStatusService struct {
	services.UnimplementedStatusServiceServer
}

// NewStatusService - Creates a new status service
func NewStatusService() services.StatusServiceServer {
	return &ovsStatusService{}

}

// OverseerStatus - tests connection
func (srv *ovsStatusService) OverseerStatus(ctx context.Context, msg *empty.Empty) (*services.ActionResultMsg, error) {

	result := &services.ActionResultMsg{
		Success: true, Message: "ovs status service ready",
	}

	return result, nil
}

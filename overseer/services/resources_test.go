package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/internal/resources"
	"github.com/przebro/overseer/proto/services"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc/codes"
)

type mockResourceManager struct {
	mock.Mock
}

func (m *mockResourceManager) Add(ctx context.Context, name string, odate date.Odate) (bool, error) {
	args := m.Called(ctx, name, odate)

	return args.Bool(0), args.Error(1)
}
func (m *mockResourceManager) Delete(ctx context.Context, name string, odate date.Odate) (bool, error) {
	args := m.Called(ctx, name, odate)
	return args.Bool(0), args.Error(1)
}
func (m *mockResourceManager) Set(ctx context.Context, name string, policy uint8) (bool, error) {
	args := m.Called(ctx, name, policy)
	return args.Bool(0), args.Error(1)
}
func (m *mockResourceManager) DestroyFlag(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)

}
func (m *mockResourceManager) ListResources(ctx context.Context, filter resources.ResourceFilter) []resources.ResourceModel {

	args := m.Called(ctx, filter)
	return args.Get(0).([]resources.ResourceModel)
}

type resourcesTestSuite struct {
	service services.ResourceServiceServer
	rsrvc   *ovsResourceService
	suite.Suite
	manager *mockResourceManager
}

func TestResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(resourcesTestSuite))
}

func (s *resourcesTestSuite) SetupSuite() {

	s.manager = &mockResourceManager{}
	s.service = NewResourceService(s.manager)
	s.rsrvc = s.service.(*ovsResourceService)

}

func (suite *resourcesTestSuite) TestAddTicket_Errors() {
	client := suite.service
	msg := &services.TicketActionMsg{}

	_, err := client.AddTicket(context.Background(), msg)

	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	msg.Odate = "20201120"

	_, err = client.AddTicket(context.Background(), msg)

	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	msg.Odate = "ABCDEF"
	msg.Name = "test"

	_, err = client.AddTicket(context.Background(), msg)

	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	msg.Name = "ticket_with_very_long_name_that_exceeds_32_characters"
	msg.Odate = "20201115"

	_, err = client.AddTicket(context.Background(), msg)

	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)
}

func (suite *resourcesTestSuite) TestAddTicket_Exists_Errors() {

	ctx := context.Background()
	msg := &services.TicketActionMsg{Name: "service_add_ticket_AB"}

	suite.manager.On("Add", ctx, msg.Name, date.Odate("")).Return(true, nil).Once()
	suite.manager.On("Add", ctx, msg.Name, date.Odate("")).Return(false, errors.New("ticket with given name and odate already exists")).Once()

	client := suite.service

	_, err := client.AddTicket(ctx, msg)

	suite.Nil(err)
	_, err = client.AddTicket(context.Background(), msg)

	suite.NotNil(err)

	_, code := matchExpectedStatusFromError(err, codes.FailedPrecondition)
	suite.Equal(codes.FailedPrecondition, code)

}

func (suite *resourcesTestSuite) TestAddTicket() {

	client := suite.service
	ctx := context.Background()
	msg := &services.TicketActionMsg{Name: "service_add_ticket_01"}
	suite.manager.On("Add", ctx, msg.Name, date.Odate("")).Return(true, nil).Once()

	r, err := client.AddTicket(ctx, msg)

	suite.Nil(err)
	suite.True(r.Success)
}

func (suite *resourcesTestSuite) TestDeleteTicket_Errors() {

	client := suite.service
	ctx := context.Background()
	msg := &services.TicketActionMsg{Name: "service_test_2", Odate: ""}

	msg.Name = "very_long_name_that_exceeds_32_characters"
	msg.Odate = ""

	_, err := client.DeleteTicket(ctx, msg)

	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	msg.Name = "service_test_2"
	msg.Odate = "ABCDEDF"

	_, err = client.DeleteTicket(ctx, msg)

	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

}

func (suite *resourcesTestSuite) TestDeleteicket_Exists_Errors() {

	client := suite.service
	ctx := context.Background()
	msg := &services.TicketActionMsg{Name: "service_test_2", Odate: ""}
	msg.Name = "service_test_2"
	msg.Odate = ""

	suite.manager.On("Delete", ctx, msg.Name, date.Odate("")).Return(true, nil).Once()
	suite.manager.On("Delete", ctx, msg.Name, date.Odate("")).Return(false, errors.New("ticket with given name and odate already exists")).Once()

	_, err := client.DeleteTicket(ctx, msg)

	suite.Nil(err)

	msg.Name = "service_test_2"
	msg.Odate = ""

	_, err = client.DeleteTicket(ctx, msg)

	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.FailedPrecondition)
	suite.Equal(codes.FailedPrecondition, code)

}

func (suite *resourcesTestSuite) TestDeleteTicket() {

	client := suite.service
	ctx := context.Background()
	msg := &services.TicketActionMsg{Name: "service_test_DEL_03", Odate: "20201120"}
	suite.manager.On("Delete", ctx, msg.Name, date.Odate(msg.Odate)).Return(true, nil).Once()

	r, err := client.DeleteTicket(ctx, msg)
	suite.Nil(err)

	suite.True(r.Success)

}

/*
func TestSetFlag_Errors(t *testing.T) {

		client := createResourceClient(t)
		msg := &services.FlagActionMsg{Name: "very_long_resource_name_that_exceeds_32_chracters"}

		_, err := client.SetFlag(context.Background(), msg)

		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
			t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
		}

		msg.Name = "test_flag_02"
		msg.State = 2
		_, err = client.SetFlag(context.Background(), msg)

		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
			t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
		}
	}

func TestSetFlag_InvalidState(t *testing.T) {

		client := createResourceClient(t)
		msg := &services.FlagActionMsg{Name: "test_flag_002A", State: 0}

		r, err := client.SetFlag(context.Background(), msg)

		if err != nil {
			t.Error("unexpected result:", err)
		}
		if r.Success != true {
			t.Error("unexpected result:", r.Success, "expected:", true)
		}

		msg.Name = "test_flag_002A"
		msg.State = 1

		//flag is already set to shared so, setting to exclusive is not allowed
		_, err = client.SetFlag(context.Background(), msg)

		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.FailedPrecondition); !ok {
			t.Error("unexpected result:", code, "expected:", codes.FailedPrecondition)
		}
	}

func TestSetFlag(t *testing.T) {

		client := createResourceClient(t)
		msg := &services.FlagActionMsg{Name: "test_flag_002C", State: 1}

		r, err := client.SetFlag(context.Background(), msg)

		if err != nil {
			t.Error("unexpected result:", err)
		}
		if r.Success != true {
			t.Error("unexpected result:", r.Success, "expected:", true)
		}
	}

func TestListFlags_Errors(t *testing.T) {

	client := createResourceClient(t)

	msg := &services.FlagActionMsg{Name: "very_long_resource_name_that_exceeds_32_chracters"}
	result, err := client.ListFlags(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = result.Recv()
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

}
func TestListFlags(t *testing.T) {

	client := createResourceClient(t)

	msg := &services.FlagActionMsg{Name: "test_flag_05A", State: 1}
	_, err := client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	msg.Name = "test_flag*"
	result, err := client.ListFlags(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = result.Recv()
	if err != nil {
		t.Error("unexpected result:", err)
	}

}

func TestDestroyFlag_Errors(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.FlagActionMsg{Name: "very_long_resource_name_that_exceeds_32_chracters"}
	_, err := client.DestroyFlag(context.Background(), msg)

	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	msg.Name = "test_flag_1N"
	_, err = client.DestroyFlag(context.Background(), msg)

	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.NotFound); !ok {
		t.Error("unexpected result:", code, "expected:", codes.NotFound)
	}

}
func TestDestroyFlag(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.FlagActionMsg{Name: "test_flag_05", State: 1}
	_, err := client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	msg.Name = "test_flag_05"
	r, err := client.DestroyFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err, "expected:nil")
	}
	if !r.Success {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

}
*/
func (suite *resourcesTestSuite) TestAllowedActions() {

	tdata := map[string]auth.UserAction{
		"AddTicket":     auth.ActionAddTicket,
		"DeleteTicket":  auth.ActionRemoveTicket,
		"SetFlag":       auth.ActionSetFlag,
		"DestroyFlag":   auth.ActionSetFlag,
		"ListResources": auth.ActionBrowse,
	}

	for k, v := range tdata {

		act := suite.rsrvc.GetAllowedAction(k)
		if act != v {
			suite.Fail(fmt.Sprintf("unexpected result:%v expected:%v", act, v))
		}

	}
}

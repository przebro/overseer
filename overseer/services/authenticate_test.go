package services

import (
	"context"
	"net"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/proto/services"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/test/bufconn"
)

var authcl services.AuthenticateServiceClient
var asrvc *ovsAuthenticateService

func createAuthClient(t *testing.T) services.AuthenticateServiceClient {

	if authcl != nil {
		return authcl
	}

	listener := bufconn.Listen(1)
	mocksrv := &mockBuffconnServer{grpcServer: grpc.NewServer(buildUnaryChain(), buildStreamChain())}

	logger.NewTestLogger()
	var err error

	tcv, err := NewTokenCreatorVerifier(authcfg)

	if err != nil {
		t.Fatal("unable to create connection", err)
	}

	authservice, err := NewAuthenticateService(authcfg, tcv, provider, logger.NewTestLogger())

	if err != nil {
		t.Error(err)
	}
	asrvc = authservice.(*ovsAuthenticateService)

	services.RegisterAuthenticateServiceServer(mocksrv.grpcServer, authservice)

	dialer := func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		t.Fatal("unable to create connection", err)
	}

	authcl = services.NewAuthenticateServiceClient(conn)
	go mocksrv.grpcServer.Serve(listener)

	return authcl

}

func TestCreateNewAuthenticateService_Error(t *testing.T) {

	prov := &datastore.Provider{}
	_, err := NewAuthenticateService(authcfg, nil, prov, logger.NewTestLogger())
	if err == nil {
		t.Error("unexpected result:", nil, "expected: not nil")
	}

}
func TestAuthenticate_Anonymous_User_Success(t *testing.T) {
	client := createAuthClient(t)

	msg := &services.AuthorizeActionMsg{Username: "", Password: ""}
	r, err := client.Authenticate(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Message != "anonymous access" {
		t.Error("unexpected result:", r.Message, "expected:", "anonymous access")
	}

}
func TestAuthenticate_Anonymous_User_Fail(t *testing.T) {
	client := createAuthClient(t)

	msg := &services.AuthorizeActionMsg{Username: "", Password: ""}

	asrvc.allowAnonymous = false

	_, err := client.Authenticate(context.Background(), msg)
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

}
func TestAuthenticate_User_Fail(t *testing.T) {

	client := createAuthClient(t)
	asrvc.allowAnonymous = false

	msg := &services.AuthorizeActionMsg{Username: "testuser1", Password: ""}

	_, err := client.Authenticate(context.Background(), msg)
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

	msg.Username = "testuser1"
	msg.Password = "invalid_password"

	_, err = client.Authenticate(context.Background(), msg)
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}
}
func TestAuthenticate_User_Success(t *testing.T) {

	client := createAuthClient(t)
	asrvc.allowAnonymous = false

	msg := &services.AuthorizeActionMsg{Username: "testuser1", Password: "notsecure"}

	r, err := client.Authenticate(context.Background(), msg)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	if r.Message == "" {
		t.Error("unexpected result token is empty")
	}

}

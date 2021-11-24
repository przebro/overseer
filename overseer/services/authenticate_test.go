package services

import (
	"context"
	"net"
	"overseer/common/logger"
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

func TestAuthenticateUser(t *testing.T) {

	client := createAuthClient(t)

	msg := &services.AuthorizeActionMsg{Username: "", Password: ""}
	r, err := client.Authenticate(context.Background(), msg)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	asrvc.allowAnonymous = false

	_, err = client.Authenticate(context.Background(), msg)
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

	msg.Username = "testuser1"

	_, err = client.Authenticate(context.Background(), msg)
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

	msg.Username = "testuser1"
	msg.Password = "notsecure"

	r, err = client.Authenticate(context.Background(), msg)
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

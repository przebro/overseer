package services

import (
	"context"
	"net"
	"overseer/proto/services"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var cl services.StatusServiceClient

func createCLient(t *testing.T) services.StatusServiceClient {

	if cl != nil {
		return cl
	}

	listener := bufconn.Listen(1)
	mocksrv := &mockBuffconnServer{grpcServer: grpc.NewServer(buildUnaryChain(), buildStreamChain())}

	services.RegisterStatusServiceServer(mocksrv.grpcServer, NewStatusService())

	dialer := func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		t.FailNow()
	}

	cl = services.NewStatusServiceClient(conn)
	go mocksrv.grpcServer.Serve(listener)

	return cl
}
func TestStausService(t *testing.T) {

	client := createCLient(t)

	response, err := client.OverseerStatus(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if response.Success != true {
		t.Error("unexpected result:", response.Success)
	}

}

package server

import (
	"context"
	"fmt"
	"net"
	"testing"

	api "github.com/pennsieve/jit-calculation-service/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	server, err := NewGRPCServer()
	require.NoError(t, err)

	go func() {
		server.Serve(l)
	}()

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	cc, err := grpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)
	client := api.NewJitClient(cc)

	response, err := client.Calculate(ctx, &api.CalculationRequest{Input: 4})
	require.NoError(t, err)
	fmt.Println(response)

	defer teardown(server, cc, l)
}

func teardown(server *grpc.Server, cc *grpc.ClientConn, l net.Listener) {
	server.Stop()
	cc.Close()
	l.Close()
}

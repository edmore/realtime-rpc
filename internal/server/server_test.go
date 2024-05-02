package server

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	api "github.com/edmore/realtime-rpc/api/v1"
	"github.com/stretchr/testify/assert"
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

	stream, err := client.Calculate(ctx)
	if err != nil {
		log.Fatalf("Failed to setup stream : %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive result : %v", err)
			}
			expected := int64(16)
			assert.Equal(t, expected, in.Result)
		}
	}()

	if err := stream.Send(&api.CalculationRequest{Input: 4}); err != nil {
		log.Fatalf("Failed to send an input: %v", err)
	}

	stream.CloseSend()
	<-waitc

	defer teardown(server, cc, l)
}

func teardown(server *grpc.Server, cc *grpc.ClientConn, l net.Listener) {
	server.Stop()
	cc.Close()
	l.Close()
}

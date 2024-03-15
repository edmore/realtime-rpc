package server

import (
	"context"

	api "github.com/pennsieve/jit-calculation-service/api/v1"
	"google.golang.org/grpc"
)

var _ api.JitServer = (*grpcServer)(nil)

func NewGRPCServer() (*grpc.Server, error) {
	gsrv := grpc.NewServer()
	srv, err := newgrpcServer()

	if err != nil {
		return nil, err
	}

	api.RegisterJitServer(gsrv, srv)

	return gsrv, nil
}

type grpcServer struct {
	api.UnimplementedJitServer
}

func newgrpcServer() (srv *grpcServer, err error) {
	srv = &grpcServer{}
	return srv, nil
}

func (s *grpcServer) Calculate(ctx context.Context, req *api.CalculationRequest) (*api.CalculationResponse, error) {

	// delegate to library
	result := req.Input * req.Input

	return &api.CalculationResponse{Result: result}, nil
}

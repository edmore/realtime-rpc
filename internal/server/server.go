package server

import (
	"io"
	"log"

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

func (s *grpcServer) Calculate(stream api.Jit_CalculateServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		result := in.Input * in.Input
		log.Println("result", result)
		if err := stream.Send(&api.CalculationResponse{Result: result}); err != nil {
			return err
		}
	}
}

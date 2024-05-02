package main

import (
	"log"
	"net"

	"github.com/edmore/realtime-rpc/internal/server"
)

const (
	port = ":50051"
)

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("something went wrong")
	}

	server, err := server.NewGRPCServer()
	if err != nil {
		log.Fatalf("something went wrong")
	}

	log.Println("gRPC server listening on port:", l.Addr().String())
	if err := server.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

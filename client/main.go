package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	api "github.com/edmore/realtime-rpc/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address = "localhost:50051"
)

func main() {
	ctx := context.Background()

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	log.Println("dialing")
	if ep, ok := os.LookupEnv("SERVER_ENDPOINT"); ok {
		address = ep
	}
	cc, err := grpc.Dial(address, clientOptions...)
	if err != nil {
		log.Fatalf("something went wrong")
	}
	defer cc.Close()

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
			fmt.Println(in.Result)
		}
	}()

	if err := stream.Send(&api.CalculationRequest{Input: 4}); err != nil {
		log.Fatalf("Failed to send an input: %v", err)
	}

	stream.CloseSend()
	<-waitc
}

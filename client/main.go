package main

import (
	"context"
	"fmt"
	"io"
	"log"

	api "github.com/pennsieve/jit-calculation-service/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	log.Println("dialing")
	cc, err := grpc.Dial(":50051", clientOptions...)
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

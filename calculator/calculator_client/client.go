package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close() // Finally close the connection

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)
	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.CalculationRequest{
		Calculation: &calculatorpb.Calculation{
			FirstNumber:  3,
			SecondNumber: 5,
		},
	}

	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculation RPC: %v", err)
	}
	log.Printf("Response from Calculation: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v", msg)
	}
}

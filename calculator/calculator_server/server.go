package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"../calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculationRequest) (*calculatorpb.CalculationResponse, error) {
	fmt.Printf("Calculate function was invoked with %v\n", req)
	firstNumber := req.GetCalculation().GetFirstNumber()
	secondNumber := req.GetCalculation().GetSecondNumber()
	result := firstNumber + secondNumber
	res := &calculatorpb.CalculationResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)
	k := 2
	number := req.GetNumber()
	for number > 1 {
		if number%int64(k) == 0 {
			res := &calculatorpb.PrimeNumberResponse{
				Number: int64(k),
			}
			stream.Send(res)
			number = number / int64(k)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k = k + 1
			fmt.Printf("Divisor has increased to %v\n", k)
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

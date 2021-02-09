package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"../calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request %v\n", stream)
	result := float32(0)
	count := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			result = result / float32(count)
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Number: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		number := req.GetNumber()
		count++
		result += float32(number)
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request %v\n", stream)

	// temp := 0
	result := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}
		number := req.GetNumber()
		if number > result {
			result = number
			err = stream.Send(&calculatorpb.FindMaxResponse{
				Number: result,
			})
			if err != nil {
				log.Fatalf("Error while sending data to client: %v\n", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function was invoked with %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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

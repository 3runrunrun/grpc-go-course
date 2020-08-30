package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/3runrunrun/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	fmt.Printf("Sum function invoked with %v\n", req)

	numA := req.GetNumA()
	numB := req.GetNumB()

	result := numA + numB

	res := &calculatorpb.SumResponse{
		Result: result,
	}

	return res, nil
}

func (*server) PrimeDecomposition(req *calculatorpb.PrimeRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	var factor int32
	factor = 2

	fmt.Printf("PrimeDecomposition invoked with: %v\n", req)

	numA := req.GetNumA()

	for numA > 1 {
		if numA%factor == 0 {
			res := &calculatorpb.PrimeResponse{
				Result: factor,
			}
			numA = int32(numA / factor)
			// send stream response
			stream.Send(res)
		} else {
			factor++
		}
	}

	return nil
}

func main() {
	fmt.Println("calculator server is running...")

	// define a connection object
	lis, err := net.Listen("tcp", "0.0.0.0:30037")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	// create grpc server
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// serve grpc server with defined connection
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

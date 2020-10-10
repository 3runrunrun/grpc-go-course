package main

import (
	"context"
	"fmt"
	"io"
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
	factor := int32(2)

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
			log.Printf("factor number has increased to %v ", factor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage rpc is invoked with streaming-client request")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("error while receiving client-stream request message: %v\n", err)
		}
		sum += req.GetNumA()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum rpc is invoked with BiDi request")

	curNumber := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client request: %v\n", err)
			return err
		}

		if req.GetNumA() > curNumber {
			curNumber = req.GetNumA()

			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Result: curNumber,
			})
			if err != nil {
				log.Fatalf("error while sending stream to the client: %v\n", err)
				return err
			}
		}

		curNumber = req.GetNumA()
	}

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

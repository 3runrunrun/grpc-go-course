package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/3runrunrun/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {

	// define a client connection, by dialing RPC target server
	cc, err := grpc.Dial("localhost:30037", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
	}

	defer cc.Close()

	// define a client service
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("client service created: %f\n", c)

	// doSum(c)
	primeFactor(c)
}

func doSum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("sum 2 number...")

	req := &calculatorpb.SumRequest{
		NumA: 3,
		NumB: 10,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to Sum: %v\n", err)
	}

	log.Printf("result = %v\n", res.Result)
}

func primeFactor(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("decomposing prime number...")

	req := &calculatorpb.PrimeRequest{
		NumA: 120,
	}
	stream, err := c.PrimeDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get stream: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of stream
			break
		}

		if err != nil {
			log.Fatalf("failed to read streams' message: %v", err)
		}

		fmt.Printf("%v ", res.GetResult())
	}

	fmt.Println("")
}

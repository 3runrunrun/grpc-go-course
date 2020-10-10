package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// primeFactor(c)
	// computeAverage(c)
	findMax(c)
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
		NumA: 12,
	}
	stream, err := c.PrimeDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get stream: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break // we've reached the end of stream
		}
		if err != nil {
			log.Fatalf("failed to read streams' message: %v", err)
		}
		fmt.Printf("%v ", res.GetResult())
	}

	fmt.Println("")
}

func computeAverage(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("starting a client-streaming (ComputeAverage) RPC...")

	numbers := []int32{1, 2, 3, 4}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("failed to call client-stream: %v\n", err)
	}

	for _, number := range numbers {
		fmt.Printf("sending number = %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			NumA: number,
		})
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive server-reponse: %v\n", err)
	}

	fmt.Printf("Average result is %v\n", res)
}

func findMax(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting BiDi streaming client...")

	inputs := []int32{1, 5, 3, 6, 2, 20}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while creating client stream: %v\n", err)
	}

	waitc := make(chan struct{})

	// ** send bunch of messages to the client
	go func() {
		// ** send bunch of messages (to the server)
		for _, input := range inputs {
			fmt.Println("---------------------------------------")
			fmt.Printf("sending number: %d\n", input)

			// create request message
			req := &calculatorpb.FindMaximumRequest{
				NumA: input,
			}
			// send message
			stream.Send(req)

			time.Sleep(1000 * time.Millisecond)
		}

		// close the stream
		stream.CloseSend()
	}()

	//  ** receive bunch of messages from the server
	go func() {
		// ** receive bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving response from the server: %v\n", err)
			}

			fmt.Printf("highest number is: %d\n", res.Result)
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/3runrunrun/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("hello, I'm a client")

	// cc >> connection client, a connection object to dial or import grpc server
	// param 01 >> grpc server target
	// param 02 >> options, WithInsecure means we dial a server without SSL. Remove it in production matter
	cc, err := grpc.Dial("localhost:30036", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	// client, a grpc client
	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("created client: %f\n", c)

	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do an Unary RPC...")

	// request message
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Fathir",
			LastName:  "Qisthi",
		},
	}

	// invoke Greet RPC
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("response from Greet RPC: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting server Streaming RPC...")

	// create stream request message
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "fathir",
			LastName:  "qisthi",
		},
	}

	// call RPC stream function
	// streamClient = RPC streams' response
	streamClient, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to invoke stream RPC: %v\n", err)
	}

	// read streams' message
	for {
		// call Receive function, to read streams' message
		msg, err := streamClient.Recv()
		// end of file OR end of stream msg
		if err == io.EOF {
			// we have reached end of stream
			break
		}

		// if any random error appear
		if err != nil {
			log.Fatalf("failed to read stream: %v\n", err)
		}

		// read streams' response
		fmt.Printf("stream response: %v\n", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a client streaming RPC...")

	// request message
	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "fathir",
				LastName:  "qisthi",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "carl",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "johnson",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "billy",
			},
		},
	}

	// invoke LongGreet RPC client streamer
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet client streamer: %v\n", err)
	}
	for _, req := range reqs { // send request streams, one by one
		fmt.Printf("sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	// receive LongGreet response and close client invoker
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving LongGreet response: %v\n", err)
	}
	fmt.Printf("LongGreet responses: %v\n", res)
}

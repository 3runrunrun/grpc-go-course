package main

import (
	"context"
	"fmt"
	"log"

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

	doUnary(c)
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

}

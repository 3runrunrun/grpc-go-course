package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/3runrunrun/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

// define grpc server func, which is accept grpc request, then return grpc response
// by implementing GreetServiceServer, which contains Greet()
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("Greet function is invoked with %v\n", req) // will be appeared when the client invokes this function

	firstname := req.GetGreeting().GetFirstName() // extract firstname from input request
	result := "Hello" + firstname

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("hello world")

	// define connection
	lis, err := net.Listen("tcp", "0.0.0.0:30036")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// define grpc server
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	// serve grpc server with defined connection
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

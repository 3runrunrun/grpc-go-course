package main

import (
	"fmt"
	"log"
	"net"

	"github.com/3runrunrun/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

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

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

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

// implementing GreetManyTimes rpc Stream
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Printf("GreetManyTimes is invoked with %v\n", req)

	firstname := req.Greeting.GetFirstName()

	for i := 0; i < 10; i++ {

		result := "hello " + firstname + " number " + strconv.Itoa(i)

		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

// implementing LongGreet rpc, the stream-client rpc
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Printf("LongGreet rpc is invoked with the streaming-client request\n")

	result := "Hello "
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			}) // we've finished reading the client-stream request
		}
		if err != nil {
			log.Fatalf("error while reading the client-stream: %v\n", err)
			return err
		}

		firstname := req.GetGreeting().GetFirstName()
		result += firstname + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone rpc is invoked with the bidi-streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"

		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalf("error while sending data to client: %v", err)
			return nil
		}
	}
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

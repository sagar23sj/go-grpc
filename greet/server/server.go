package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/sagar23sj/go-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (resp *greetpb.GreetResponse, err error) {
	fmt.Println("Greet Function was invoked to demonstrate unary streaming")

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	res := "Hello " + firstName + lastName

	resp = &greetpb.GreetResponse{
		Result: res,
	}
	return resp, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, resp greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Println("GreetManyTimes function invoked for server side streaming")

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		result := fmt.Sprintf("Hey There!, My name is %v %v and my id is #%v", firstName, lastName, i)
		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}

		time.Sleep(1000 * time.Millisecond)
		resp.Send(&res)

	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Println("LongGreet Function is invoked to demonstrate client side streaming")

	result := ""

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}

		firstName := msg.GetGreeting().GetFirstName()
		lastName := msg.GetGreeting().GetLastName()

		result += "\nHello " + firstName + " " + lastName + " !"
	}
	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone Function is invoked to demonstrate Bi-directional streaming\n")

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving data from GreetEveryone client : %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()

		result := "Hello " + firstName + " " + lastName + " !"

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if sendErr != nil {
			log.Fatalf("error while sending response to GreetEveryone Client : %v", err)
			return err
		}
	}
	return nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {

	fmt.Println("\nInvoking GreetWithDeadline Function to demonstrate deadlines in GRPC")

	result := ""

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result = "Hello " + firstName + " " + lastName + " !"
	time.Sleep(3 * time.Second)
	return &greetpb.GreetWithDeadlineResponse{Result: result}, nil
}

func main() {
	fmt.Println("vim-go")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatal("failed to serve : %v", err)
	}
}

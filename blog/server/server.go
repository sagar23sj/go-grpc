package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"githu.com/sagar23sj/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Initiating Blog gRPC Service.......")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {

		fmt.Println("Starting Blog Server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to start service : %v", err)
		}
	}()

	//wait for Ctrl+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block Until a Signal is Received
	<-ch
	fmt.Println("Stopping Blog Service Server.....")
	s.Stop()

	fmt.Println("Closing the Listener.....")
	lis.Close()
}

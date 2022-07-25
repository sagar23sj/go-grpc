package main

import (
	"log"

	"google.golang.org/grpc"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatatlf("could not connect to server : %v", err)
	}
	defer cc.Close()

}

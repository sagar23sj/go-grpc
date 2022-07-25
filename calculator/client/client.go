package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sagar23sj/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	conn := calculatorpb.NewCalculatorServiceClient(cc)

	//Unary streaming operation - Sum
	sumNumbers(conn)

	//Server Streaming operation - PrimeFactorDecomposition
	primefactorDecomposition(conn)

	//Client Streaming Opeartion - ComputeAverage
	computeAverage(conn)

	//Bi-Directional Streaming Operation - FindMaxNumber
	findMaxNumber(conn)

	//Error Handling in GRPC demonstration - SquareRoot
	squareRoot(conn)
}

func sumNumbers(conn calculatorpb.CalculatorServiceClient) {
	fmt.Println("\n**********************************************")
	fmt.Println("Initiating Sum of 2 Numbers using unary GRPC")

	firstNo := 3
	secondNo := 4
	req := calculatorpb.SumRequest{
		FirstNumber:  int32(firstNo),
		SecondNumber: int32(secondNo),
	}

	resp, err := conn.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling sum function over grpc server : %v", err)
	}

	fmt.Printf("Sum of %v and %v = %v", firstNo, secondNo, resp.SumResult)
}

func primefactorDecomposition(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("\n**********************************************")
	fmt.Println("Initiating Prime Number decompostion using server streaming")

	number := int32(95)

	req := calculatorpb.PrimeNumberDecompositionRequest{
		Number: number,
	}

	resp, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling primefactorDecomposition ovr grpc server : %v", err)
	}

	fmt.Println("Prime Factor Decompostion results :")
	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("error occured while reading response from server : %v", err)
		}
		fmt.Print(msg.GetPrimeFactor())
	}
}

func computeAverage(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("\n**********************************************")
	fmt.Println("Initiating Client Side Streaming to calculate average of sent numbers in reqeust ")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error calling compute average function : %v", err)
	}

	count := 0
	for i := 100; count < 10; i = i - 3 {
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: int32(i),
		})

		if err != nil {
			log.Fatalf("error while sending request to server : %v", err)
		}

		count++
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response from ComputeAverage Server : %v", err)
	}

	fmt.Println("Average of Numbers is : ", response.GetAverage())

}

func findMaxNumber(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("\n**********************************************")
	fmt.Println("Initiating FindMaxNumber using Bi-directional Streaming ")

	stream, err := c.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("error calling findMaxNumber function : %v", err)
	}

	requests := []*calculatorpb.FindMaxNumberRequest{
		&calculatorpb.FindMaxNumberRequest{
			Number: 3,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 23,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 21,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 12,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 34,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 45,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 33,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 30,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 67,
		},
		&calculatorpb.FindMaxNumberRequest{
			Number: 80,
		},
	}

	//go rountine to send numbers
	go func([]*calculatorpb.FindMaxNumberRequest) {

		for _, req := range requests {

			fmt.Println("Sending Request to server .....")
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending findMaxNumber request to server : %v", err)
			}

			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	waitchan := make(chan struct{})

	//go routine to receive max number
	go func() {

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("error while receiving response from server : %v", err)
			}

			fmt.Println("Current Maximum : ", resp.GetMaxNumber())
		}
	}()

	<-waitchan
}

func squareRoot(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("\n**********************************************")
	fmt.Println("Initiating SquareRoot Function to demonstrate Error Handline in GRPC")

	doSquareRoot(c, 1000)
	doSquareRoot(c, -9)
}

func doSquareRoot(c calculatorpb.CalculatorServiceClient, number int64) {

	resp, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: number,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from GRPC (user error)
			fmt.Println("\nError Message : ", respErr.Message())
			fmt.Println("\nError Code : ", respErr.Code())

			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("big error calling sqaureroot : %v", err)
			return
		}
	}

	fmt.Printf("SquareRoot Result : %v", resp.GetNumberRoot())
}

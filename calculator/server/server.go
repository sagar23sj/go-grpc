package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/sagar23sj/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (resp *calculatorpb.SumResponse, err error) {
	fmt.Println("=>Executing Sum function via GRPC ......")

	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()

	result := firstNumber + secondNumber

	resp = &calculatorpb.SumResponse{
		SumResult: result,
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, resp calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	fmt.Println("=>Executing Server Streaming of Prime Number Decompostion over GRPC....")

	number := req.GetNumber()
	divisor := int32(2)
	for number > 1 {
		if number%divisor == 0 {
			resp.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	fmt.Println("=>Executing ComputeAverage on Server....")

	counter := 0
	sum := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//when all the requests are received

			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float32(sum / counter),
			})
		}

		if err != nil {
			log.Fatalf("error receiving requests from the client : %v", err)
		}

		sum += int(req.GetNumber())
		counter++
	}

	return nil
}

func (*server) FindMaxNumber(stream calculatorpb.CalculatorService_FindMaxNumberServer) error {

	fmt.Println("=>Executing FindMaxNumber on Server implementing Bi-Di streaming")

	maxNumber := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error receiving request data from client side : %v", err)
			return err
		}

		number := req.GetNumber()

		if number > maxNumber {
			maxNumber = number
		}

		sendErr := stream.Send(&calculatorpb.FindMaxNumberResponse{
			MaxNumber: int32(maxNumber),
		})

		if sendErr != nil {
			log.Fatalf("error sending data to FindMaxNumber client : %v", err)
			return sendErr
		}

	}

	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	fmt.Println("=>Executing SquareRoot Function on Sever implementing Error Handling")

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number : %v", number),
		)

	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Starting GRPC Server..........")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	//Register reflection service on gRPC server
	reflection.Register(s)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Serve : %v", err)
	}
}

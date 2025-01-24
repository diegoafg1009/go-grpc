package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	calculatorv1 "github.com/diegoafg1009/go-grpc/proto/generated/calculator/v1"
	"github.com/diegoafg1009/go-grpc/proto/generated/calculator/v1/calculatorv1connect"
)

const addr = "http://0.0.0.0:8080"

func main() {
	client := calculatorv1connect.NewCalculatorServiceClient(http.DefaultClient, addr)

	request := connect.NewRequest(&calculatorv1.SumRequest{
		Numbers: []int32{1, 2, 3, 4, 5},
	})

	response, err := client.Sum(context.Background(), request)

	if err != nil {
		log.Fatalf("Failed to sum: %v", err)
		return
	}

	fmt.Println("Sum:", response.Msg.Result)

	stream, err := client.PrimeNumberDecomposition(context.Background(), connect.NewRequest(&calculatorv1.PrimeNumberDecompositionRequest{
		Number: 120,
	}))

	if err != nil {
		log.Fatalf("Failed to prime number decomposition: %v", err)
		return
	}

	for stream.Receive() {
		fmt.Println("PrimeNumberDecomposition:", stream.Msg().PrimeFactor)
	}
}

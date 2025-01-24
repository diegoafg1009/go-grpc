package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	calculatorv1 "github.com/diegoafg1009/go-grpc/proto/generated/calculator/v1"
	"github.com/diegoafg1009/go-grpc/proto/generated/calculator/v1/calculatorv1connect"
	greetv1 "github.com/diegoafg1009/go-grpc/proto/generated/greet/v1"
	"github.com/diegoafg1009/go-grpc/proto/generated/greet/v1/greetv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const addr = "localhost:8080"

func main() {
	mux := http.NewServeMux()

	path, handler := greetv1connect.NewGreetServiceHandler(GreetService{})
	mux.Handle(path, handler)

	path, handler = calculatorv1connect.NewCalculatorServiceHandler(CalculatorService{})
	mux.Handle(path, handler)

	fmt.Println("... Listening on", addr)
	http.ListenAndServe(
		addr,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

type GreetService struct {
	greetv1connect.UnimplementedGreetServiceHandler
}

func (GreetService) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	firstName := req.Msg.FirstName
	greeting := fmt.Sprintf("Hello, %s!", firstName)
	return connect.NewResponse(&greetv1.GreetResponse{
		Greeting: greeting,
	}), nil
}

func (GreetService) GreetManyTimes(ctx context.Context, req *connect.Request[greetv1.GreetRequest], stream *connect.ServerStream[greetv1.GreetResponse]) error {
	firstName := req.Msg.FirstName
	for i := 0; i < 10; i++ {
		greeting := fmt.Sprintf("Hello, %s! - %d", firstName, i)
		stream.Send(&greetv1.GreetResponse{
			Greeting: greeting,
		})
	}
	return nil
}

func (GreetService) LongGreet(ctx context.Context, stream *connect.ClientStream[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	response := ""
	for stream.Receive() {
		name := stream.Msg().FirstName
		greeting := fmt.Sprintf("Hello, %s\n", name)
		response += greeting
	}
	return connect.NewResponse(&greetv1.GreetResponse{
		Greeting: response,
	}), nil
}

type CalculatorService struct {
	calculatorv1connect.UnimplementedCalculatorServiceHandler
}

func (CalculatorService) Sum(ctx context.Context, req *connect.Request[calculatorv1.SumRequest]) (*connect.Response[calculatorv1.SumResponse], error) {
	numbers := req.Msg.Numbers
	var sum int32
	for _, number := range numbers {
		sum += number
	}
	return connect.NewResponse(&calculatorv1.SumResponse{Result: sum}), nil
}

func (CalculatorService) PrimeNumberDecomposition(ctx context.Context, req *connect.Request[calculatorv1.PrimeNumberDecompositionRequest], stream *connect.ServerStream[calculatorv1.PrimeNumberDecompositionResponse]) error {
	number := req.Msg.Number
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorv1.PrimeNumberDecompositionResponse{PrimeFactor: divisor})
			number = number / divisor
		} else {
			divisor++
		}
	}
	return nil
}

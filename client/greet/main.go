package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	greetv1 "github.com/diegoafg1009/go-grpc/proto/generated/greet/v1"
	"github.com/diegoafg1009/go-grpc/proto/generated/greet/v1/greetv1connect"
)

const addr = "http://0.0.0.0:8080"

func main() {

	client := greetv1connect.NewGreetServiceClient(http.DefaultClient, addr)

	request := connect.NewRequest(&greetv1.GreetRequest{
		FirstName: "Diego",
	})

	response, err := client.Greet(context.Background(), request)

	if err != nil {
		log.Fatalf("Failed to greet: %v", err)
		return
	}

	fmt.Println("Greet:", response.Msg.Greeting)

	stream, err := client.GreetManyTimes(context.Background(), request)

	if err != nil {
		log.Fatalf("Failed to greet many times: %v", err)
		return
	}

	for stream.Receive() {
		fmt.Println("GreetManyTimes:", stream.Msg().Greeting)
	}

	longGreetClient := client.LongGreet(context.Background())

	longGreetClient.Send(&greetv1.GreetRequest{
		FirstName: "Diego",
	})

	longGreetClient.Send(&greetv1.GreetRequest{
		FirstName: "Andre",
	})

	response, err = longGreetClient.CloseAndReceive()

	if err != nil {
		log.Fatalf("Failed to long greet: %v", err)
		return
	}

	fmt.Println("LongGreet:")
	fmt.Println(response.Msg.Greeting)
}

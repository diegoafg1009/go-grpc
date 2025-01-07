package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
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

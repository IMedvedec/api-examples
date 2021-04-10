package greeting

import (
	"context"

	"github.com/imedvedec/api-examples/api/grpc/build"
)

type greetingServer struct {
	build.UnimplementedGreetingServer
}

func New() build.GreetingServer {
	return &greetingServer{}
}

func (gs *greetingServer) Greet(ctx context.Context, request *build.GreetingRequest) (*build.GreetingResponse, error) {
	return &build.GreetingResponse{
		Message: "Ok, this works!",
	}, nil
}

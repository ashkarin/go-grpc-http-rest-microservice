package samples_storage

import (
	"context"

	api "github.com/ashkarin/go-grpc-http-rest-microservice/internal/services/samples_storage/api"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterService registers Service on GRPC Server
func RegisterService(s *grpc.Server, srv api.SamplesStorageServiceServer) {
	api.RegisterSamplesStorageServiceServer(s, srv)
}

// RegisterGateway
func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, opts []grpc.DialOption) error {
	return api.RegisterSamplesStorageServiceHandlerFromEndpoint(ctx, mux, address, opts)
}

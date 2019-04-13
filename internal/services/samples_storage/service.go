package samples_storage

import (
	"context"

	api "github.com/ashkarin/go-grpc-http-rest-microservice/internal/services/samples_storage/api"
	"github.com/ashkarin/go-grpc-http-rest-microservice/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type ServiceServer struct {
	storeSampleUsecase *usecases.StoreSampleUsecase
}

// NewServiceServer creates api.SamplesStorageServiceServer
func NewServiceServer(ucase *usecases.StoreSampleUsecase) api.SamplesStorageServiceServer {
	return &ServiceServer{
		storeSampleUsecase: ucase,
	}
}

func (s *ServiceServer) Store(ctx context.Context, req *api.StoreRequest) (*api.StoreResponse, error) {
	counter := len(req.Samples)
	var samples []interface{}
	for _, sample := range req.Samples {
		// Decode received Sample
		if sample, err := DecodeSample(sample); err == nil {
			samples = append(samples, sample)
			continue
		} else {
			// In case of any error
			log.Error(err)
			counter--
		}
	}

	// Unpack samples
	if err := s.storeSampleUsecase.Store(ctx, samples...); err != nil {
		return nil, err
	}

	return &api.StoreResponse{Number: int64(counter)}, nil
}

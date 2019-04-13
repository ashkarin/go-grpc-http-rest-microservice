package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/acceleration"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/activity"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/location"
)

type StoreSampleUsecase struct {
	locationStorage     location.StorageGateway
	activityStorage     activity.StorageGateway
	accelerationStorage acceleration.StorageGateway
}

// NewStoreSampleUsecase creates a use case of storing samples
func NewStoreSampleUsecase(
	locationStorage location.StorageGateway,
	activityStorage activity.StorageGateway,
	accelerationStorage acceleration.StorageGateway,
) *StoreSampleUsecase {
	return &StoreSampleUsecase{
		locationStorage:     locationStorage,
		activityStorage:     activityStorage,
		accelerationStorage: accelerationStorage,
	}
}

func (s *StoreSampleUsecase) Store(ctx context.Context, samples ...interface{}) error {
	var (
		errorMesage   string
		locations     []*location.Location
		activities    []*activity.Activity
		accelerations []*acceleration.Acceleration
		unknown       []interface{}
	)

	for _, sample := range samples {
		switch sample := sample.(type) {
		case *location.Location:
			locations = append(locations, sample)
		case *activity.Activity:
			activities = append(activities, sample)
		case *acceleration.Acceleration:
			accelerations = append(accelerations, sample)
		case []interface{}:
			log.Fatal("Probably, the 'samples' array was no unwrapped before passing to the usecase.")
		default:
			// This case should never happen if properly use.
			unknown = append(unknown, sample)
		}
	}

	// Probably, it is better to do a transaction here
	if err := s.locationStorage.Store(ctx, locations...); err != nil {
		return err
	}

	if err := s.activityStorage.Store(ctx, activities...); err != nil {
		return err
	}

	if err := s.accelerationStorage.Store(ctx, accelerations...); err != nil {
		return err
	}

	if len(unknown) != 0 {
		for _, sample := range unknown {
			errorMesage += fmt.Sprintf("Sample '%v' has unknown type.\n", sample)
		}
		return fmt.Errorf(errorMesage)
	}
	return nil
}

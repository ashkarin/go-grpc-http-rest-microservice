package usecases

import (
	"database/sql"

	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/acceleration"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/activity"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/location"
)

func CreateSamplesStorageUsecase(db *sql.DB) (*StoreSampleUsecase, error) {
	// Create storages gateways
	locationsStorage, err := location.NewSQLStorageGateway(db)
	if err != nil {
		return nil, err
	}
	activityStorage, err := activity.NewSQLStorageGateway(db)
	if err != nil {
		return nil, err
	}
	accelerationStorage, err := acceleration.NewSQLStorageGateway(db)
	if err != nil {
		return nil, err
	}

	// Create usecase
	sampleStoreUsecase := NewStoreSampleUsecase(locationsStorage, activityStorage, accelerationStorage)
	return sampleStoreUsecase, nil
}

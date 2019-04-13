package samples_storage

import (
	"fmt"

	api "github.com/ashkarin/go-grpc-http-rest-microservice/internal/services/samples_storage/api"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/acceleration"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/activity"
	"github.com/ashkarin/go-grpc-http-rest-microservice/pkg/location"

	"github.com/golang/protobuf/ptypes"
)

// DecodeSample tries to convert api.Sample to any of known entities
func DecodeSample(sample *api.Sample) (interface{}, error) {
	if obj, err := sampleToActivity(sample); err == nil {
		return obj, nil
	}
	if obj, err := sampleToLocation(sample); err == nil {
		return obj, nil
	}
	if obj, err := sampleToAcceleration(sample); err == nil {
		return obj, nil
	}
	return nil, fmt.Errorf("Unknown sample type: %v", sample)
}

// sampleToActivity tries to convert the api.Sample message to the models.Activity
func sampleToActivity(sample *api.Sample) (*activity.Activity, error) {
	// Unmarshal the Data filed of the api.Sample message to the api.Activity message
	data := &api.Activity{}
	if err := ptypes.UnmarshalAny(sample.Data, data); err != nil {
		return nil, err
	}

	// Convert Timestamp
	timestamp, err := ptypes.Timestamp(sample.Timestamp)
	if err != nil {
		return nil, err
	}

	// Construct the Activity
	obj := &activity.Activity{
		ID:         sample.Id,
		Timestamp:  timestamp,
		Unknown:    data.Unknown,
		Stationary: data.Stationary,
		Walking:    data.Walking,
		Running:    data.Running,
	}

	return obj, nil
}

// sampleToLocation tries to convert the api.Sample message to the models.Location
func sampleToLocation(sample *api.Sample) (*location.Location, error) {
	// Unmarshal the Data filed of the Sample message to the Location message
	data := &api.Location{}
	if err := ptypes.UnmarshalAny(sample.Data, data); err != nil {
		return nil, err
	}

	// Convert Timestamp
	timestamp, err := ptypes.Timestamp(sample.Timestamp)
	if err != nil {
		return nil, err
	}

	// Construct the Location
	obj := &location.Location{
		ID:        sample.Id,
		Timestamp: timestamp,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
	}

	return obj, nil
}

// sampleToAcceleration tries to convert the api.Sample message to the models.Acceleration
func sampleToAcceleration(sample *api.Sample) (*acceleration.Acceleration, error) {
	// Unmarshal the Data filed of the api.Sample message to the api.Acceleration message
	data := &api.Acceleration{}
	if err := ptypes.UnmarshalAny(sample.Data, data); err != nil {
		return nil, err
	}

	// Convert Timestamp
	timestamp, err := ptypes.Timestamp(sample.Timestamp)
	if err != nil {
		return nil, err
	}

	// Construct the Acceleration
	obj := &acceleration.Acceleration{
		ID:        sample.Id,
		Timestamp: timestamp,
		X:         data.X,
		Y:         data.Y,
		Z:         data.Z,
	}

	return obj, nil
}

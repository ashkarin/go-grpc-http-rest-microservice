package location

import "time"

// Location is a value acquired by geolocation device.
type Location struct {
	ID        int64     `json:"id" sql:"id"`
	Timestamp time.Time `json:"timestamp" sql:"timestamp" validate:"required"`
	Latitude  float64   `json:"latitude" sql:"latitude" validate:"required"`
	Longitude float64   `json:"longitude" sql:"longitude" validate:"required"`
}

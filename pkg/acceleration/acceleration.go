package acceleration

import "time"

// Acceleration is a value acquired by an acceleration sensor.
type Acceleration struct {
	ID        int64     `json:"id" sql:"id"`
	Timestamp time.Time `json:"timestamp" sql:"timestamp" validate:"required"`
	X         float64   `json:"x" sql:"x" validate:"required"`
	Y         float64   `json:"y" sql:"y" validate:"required"`
	Z         float64   `json:"z" sql:"z" validate:"required"`
}

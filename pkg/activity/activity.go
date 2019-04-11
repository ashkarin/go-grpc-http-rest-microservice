package activity

import "time"

// Activity is a value characterizing a person activity.
// More than one or none activity at the same time is possible.
type Activity struct {
	ID         int64     `json:"id" sql:"id"`
	Timestamp  time.Time `json:"timestamp" sql:"timestamp" validate:"required"`
	Unknown    bool      `json:"unknown" sql:"unknown"`
	Stationary bool      `json:"stationary" sql:"stationary"`
	Walking    bool      `json:"walking" sql:"walking"`
	Running    bool      `json:"running" sql:"running"`
}

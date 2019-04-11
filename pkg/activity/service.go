package activity

import "context"

// Service represents a usecase for Activity
type Service interface {
	Store(ctx context.Context, v *Activity) error
}

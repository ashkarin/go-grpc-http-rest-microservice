package acceleration

import "context"

// Service represents a usecase for Acceleration
type Service interface {
	Store(ctx context.Context, v *Acceleration) error
}

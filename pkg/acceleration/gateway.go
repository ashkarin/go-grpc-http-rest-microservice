package acceleration

import "context"

// StorageGateway represent a data storage service
type StorageGateway interface {
	Store(ctx context.Context, a ...*Acceleration) error
}

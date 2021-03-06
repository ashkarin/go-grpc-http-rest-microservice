package activity

import "context"

// StorageGateway represent a data storage service
type StorageGateway interface {
	Store(ctx context.Context, entries ...*Activity) error
}

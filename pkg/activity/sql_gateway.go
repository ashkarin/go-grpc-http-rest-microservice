package activity

import (
	"context"
	"database/sql"
)

const (
	createTableStatement = `CREATE TABLE IF NOT EXISTS ` +
		`activity("id" SERIAL PRIMARY KEY, "timestamp" TIMESTAMP, ` +
		`"unknown" BOOLEAN, "stationary" BOOLEAN, "walking" BOOLEAN, "running" BOOLEAN);`
	storeStatement = "INSERT INTO activity(id, timestamp, unknown, walking, stationary, running) VALUES($1,$2,$3,$4,$5,$6)"
)

type sqlStorageGateway struct {
	db        *sql.DB
	tableName string
}

// NewSQLStorageGateway creates a gateway to SQL Database
func NewSQLStorageGateway(db *sql.DB) (StorageGateway, error) {
	if _, err := db.Exec(createTableStatement); err != nil {
		return nil, err
	}
	return &sqlStorageGateway{db: db}, nil
}

func (r *sqlStorageGateway) Store(ctx context.Context, entries ...*Activity) error {
	stmt, err := r.db.PrepareContext(ctx, storeStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err := stmt.ExecContext(ctx, entry.ID, entry.Timestamp, entry.Unknown, entry.Stationary, entry.Walking, entry.Running)
		if err != nil {
			return err
		}
	}
	return nil
}

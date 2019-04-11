package location

import (
	"context"
	"database/sql"
)

const (
	createTableStatement = `CREATE TABLE IF NOT EXISTS ` +
		`location("id" SERIAL PRIMARY KEY, "timestamp" TIMESTAMP, ` +
		`"latitude" DOUBLE PRECISION, "longitude" DOUBLE PRECISION);`
	storeStatement = "INSERT INTO location(id, timestamp, latitude, longitude) VALUES($1,$2,$3,$4)"
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

func (r *sqlStorageGateway) Store(ctx context.Context, entries ...*Location) error {
	stmt, err := r.db.PrepareContext(ctx, storeStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err := stmt.ExecContext(ctx, entry.ID, entry.Timestamp, entry.Latitude, entry.Longitude)
		if err != nil {
			return err
		}
	}
	return nil
}

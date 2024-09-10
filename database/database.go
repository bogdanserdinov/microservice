package database

import (
	"context"
	"database/sql"

	"microservice/service"
)

var _ service.DB = (*DB)(nil)

type DB struct {
	inner *sql.DB
}

func New(inner *sql.DB) *DB {
	return &DB{
		inner: inner,
	}
}

func (db *DB) CreateDummy(ctx context.Context, dummy service.Dummy) error {
	query := `INSERT INTO dummies(id, status, description, created_at)
	          VALUES ($1, $2, $3, $4)`

	_, err := db.inner.ExecContext(ctx, query, dummy.ID, dummy.Status, dummy.Description, dummy.CreatedAt)
	return err
}

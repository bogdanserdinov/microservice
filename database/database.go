package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"microservice/service"
)

var _ service.DB = (*DB)(nil)

type DB struct {
	inner *pgxpool.Pool
}

func New(inner *pgxpool.Pool) *DB {
	return &DB{
		inner: inner,
	}
}

func (db *DB) CreateDummy(ctx context.Context, dummy service.Dummy) error {
	query := `INSERT INTO dummies(id, status, description, created_at)
	          VALUES ($1, $2, $3, $4)`

	_, err := db.inner.Exec(ctx, query, dummy.ID, dummy.Status, dummy.Description, dummy.CreatedAt)
	return err
}

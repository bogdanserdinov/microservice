package tests

import (
	"context"

	"microservice/service"
)

var _ service.DB = (*MockDB)(nil)

type MockDB struct{}

func (m MockDB) CreateDummy(ctx context.Context, dummy service.Dummy) error {
	return nil
}

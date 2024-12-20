package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"microservice/service"
)

func (s *MicroserviceSuite) TestService() {
	t := s.T()

	// Service tests covers both the service and the database.
	s.WithRealDB()

	ctx := context.Background()

	t.Run("create with invalid status", func(t *testing.T) {
		err := s.service.Create(ctx, "test", "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input value for enum status")
	})

	t.Run("create", func(t *testing.T) {
		err := s.service.Create(ctx, service.StatusSuccess, "test")
		require.NoError(t, err)
	})
}

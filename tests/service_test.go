package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"microservice/service"
)

// TestMicroserviceSuite is root function that under the hood runs all tests within the suite.
func TestMicroserviceSuite(t *testing.T) {
	suite.Run(t, new(MicroserviceSuite))
}

func (s *MicroserviceSuite) TestMicroserviceSuite() {
	t := s.T()

	ctx := context.Background()

	t.Run("create with invalid status", func(t *testing.T) {
		t.Parallel()

		err := s.service.Create(ctx, "test", "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid input value for enum status")
	})

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		err := s.service.Create(ctx, service.StatusSuccess, "test")
		require.NoError(t, err)
	})
}

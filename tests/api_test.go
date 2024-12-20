package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api_errors "microservice/pkg/http/errors"
	"microservice/servers/public/controllers"
	"microservice/service"
)

func (s *MicroserviceSuite) TestHandler() {
	t := s.T()

	// API testing is done by using mock database to check only the API functionality.
	s.WithMockDB()

	t.Run("invalid request", func(t *testing.T) {
		req := controllers.CreateRequest{
			Status:      "invalid",
			Description: "invalid",
		}

		resp := s.MakeRequest(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "invalid status", getErrorMsg(t, resp))
	})

	t.Run("successful request", func(t *testing.T) {
		req := controllers.CreateRequest{
			Status:      service.StatusSuccess,
			Description: "description",
		}

		resp := s.MakeRequest(req)
		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func (s *MicroserviceSuite) MakeRequest(reqPayload controllers.CreateRequest) *http.Response {
	t := s.T()
	t.Helper()

	rawPayload, err := json.Marshal(reqPayload)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/dummy", bytes.NewReader(rawPayload))
	recorder := httptest.NewRecorder()

	s.controller.Create(recorder, request)

	return recorder.Result()
}

func getErrorMsg(t *testing.T, resp *http.Response) string {
	t.Helper()

	var errorResp api_errors.Response
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)

	require.NoError(t, resp.Body.Close())

	return errorResp.Error
}

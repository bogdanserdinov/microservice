package controllers

import (
	"encoding/json"
	"errors"
	http_errors "microservice/pkg/http/errors"

	"net/http"

	"go.uber.org/zap"

	"microservice/service"
)

type Dummy struct {
	log *zap.Logger

	service *service.Service
}

func NewDummy(log *zap.Logger, service *service.Service) *Dummy {
	return &Dummy{
		log:     log,
		service: service,
	}
}

func (controller *Dummy) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := struct {
		Status      service.Status `json:"status"`
		Description string         `json:"description"`
	}{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http_errors.ServeError(w, http.StatusBadRequest, err)
		return
	}

	err := controller.service.Create(ctx, req.Status, req.Description)
	if err != nil {
		controller.log.Error("could not create dummy record", zap.Error(err))
		http_errors.ServeError(w, http.StatusInternalServerError, errors.New("could not create dummy entity"))
		return
	}
}

package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	http_errors "microservice/pkg/http/errors"
	"microservice/service"
)

type Dummy struct {
	log    *zap.Logger
	tracer trace.Tracer

	service *service.Service
}

func NewDummy(log *zap.Logger, tracer trace.Tracer, service *service.Service) *Dummy {
	return &Dummy{
		log:     log,
		tracer:  tracer,
		service: service,
	}
}

type CreateRequest struct {
	Status      service.Status `json:"status"`
	Description string         `json:"description"`
}

func (controller *Dummy) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := controller.tracer.Start(r.Context(), "create_http")
	defer span.End()

	req := new(CreateRequest)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil {
		http_errors.ServeError(w, http.StatusBadRequest, err)
		return
	}

	if !req.Status.IsValid() {
		http_errors.ServeError(w, http.StatusBadRequest, errors.New("invalid status"))
		return
	}

	err := controller.service.Create(ctx, req.Status, req.Description)
	if err != nil {
		controller.log.Error("could not create dummy record", zap.Error(err))
		http_errors.ServeError(w, http.StatusInternalServerError, err)
		return
	}
}

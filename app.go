package microservice

import (
	"go.uber.org/zap"
	"microservice/service"
)

type App struct {
	log *zap.Logger

	dummy *service.Service
}

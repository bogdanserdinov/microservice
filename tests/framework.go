package tests

import (
	"context"
	"database/sql"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"microservice"
	"microservice/database"
	"microservice/pkg/tracer"
	"microservice/service"
)

type MicroserviceSuite struct {
	suite.Suite

	db               *sql.DB
	tracerShutdowner func(ctx context.Context) error

	service *service.Service
}

func (s *MicroserviceSuite) SetupSuite() {
	ctx := context.Background()

	cfg := getConfig()

	db, err := sql.Open("postgres", cfg.Database.URL)
	s.Require().NoError(err)

	tracer, shutdown, err := tracer.Init(ctx, "dummy_service", cfg.Traces.JaegerEndpoint)
	if err != nil {
		log.Fatal("could not init tracer", zap.Error(err))
	}
	defer func() {
		err = shutdown(ctx)
		if err != nil {
			log.Fatal("could not shutdown the tracer", zap.Error(err))
		}
	}()

	s.db = db
	s.tracerShutdowner = shutdown
	s.service = service.New(tracer, database.New(s.db))
}

func (s *MicroserviceSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)
}

func getConfig() *microservice.Config {
	err := godotenv.Load(".test.env")
	if err != nil {
		log.Fatal("Error loading .env file", zap.Error(err))
	}

	cfg := &microservice.Config{}
	err = env.Parse(cfg)
	if err != nil {
		log.Fatal("could not parse environment variables", zap.Error(err))
	}
	return cfg
}

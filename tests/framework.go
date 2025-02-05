package tests

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"microservice"
	"microservice/database"
	"microservice/pkg/tracer"
	"microservice/servers/public/controllers"
	"microservice/service"
)

type MicroserviceSuite struct {
	suite.Suite

	cfg *microservice.Config

	db               *pgxpool.Pool
	tracer           trace.Tracer
	tracerShutdowner func(ctx context.Context) error

	service *service.Service

	controller *controllers.Dummy
}

func (s *MicroserviceSuite) SetupSuite() {
	ctx := context.Background()

	cfg := getConfig()
	s.cfg = &cfg

	db, err := pgxpool.New(ctx, s.cfg.Database.URL)
	s.Require().NoError(err)

	tracer, shutdown, err := tracer.Init(ctx, "dummy_service", s.cfg.Traces.JaegerEndpoint)
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
	s.tracer = tracer
	s.tracerShutdowner = shutdown
	s.service = service.New(tracer, database.New(s.db))
	s.controller = controllers.NewDummy(zap.NewNop(), tracer, s.service)
}

func (s *MicroserviceSuite) TearDownSuite() {
	s.db.Close()

	err := s.tracerShutdowner(context.Background())
	s.Require().NoError(err)
}

func (s *MicroserviceSuite) WithMockDB() {
	db := &MockDB{}

	s.service = service.New(s.tracer, db)
	s.controller = controllers.NewDummy(
		zap.NewNop(),
		s.tracer,
		s.service,
	)
}

func (s *MicroserviceSuite) WithRealDB() {
	db, err := pgxpool.New(context.Background(), s.cfg.Database.URL)
	s.Require().NoError(err)

	s.service = service.New(s.tracer, database.New(db))
	s.controller = controllers.NewDummy(
		zap.NewNop(),
		s.tracer,
		s.service,
	)
}

func getConfig() microservice.Config {
	err := godotenv.Load(".test.env")
	if err != nil {
		log.Fatal("error loading .env file", zap.Error(err))
	}

	cfg := &microservice.Config{}

	err = env.Parse(cfg)
	if err != nil {
		log.Fatal("could not parse environment variables", zap.Error(err))
	}

	return *cfg
}

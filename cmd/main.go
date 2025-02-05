package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"microservice"
	"microservice/pkg/tracer"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	loggerCfg := zap.NewProductionConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalln("could not build logger from config", err)
		return
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Fatalln("could not sync logger", err)
		}
	}()

	cfg, err := getConfigFromEnv()
	if err != nil {
		logger.Fatal("could not parse config", zap.Error(err))
	}

	tracer, shutdown, err := tracer.Init(ctx, "dummy_service", cfg.Traces.JaegerEndpoint)
	if err != nil {
		logger.Fatal("could not init tracer", zap.Error(err))
	}
	defer func() {
		err = shutdown(ctx)
		if err != nil {
			logger.Error("could not shutdown the tracer", zap.Error(err))
		}
	}()

	db, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		logger.Error("can't open connection to postgres", zap.Error(err))
		return
	}
	defer func() {
		db.Close()
	}()
	if err := db.Ping(ctx); err != nil {
		logger.Error("can't ping database", zap.Error(err))
		return
	}

	app := microservice.New(logger, cfg, tracer, db)

	logger.Info("servers shutdown err", zap.Error(app.Run(ctx)))
}

func getConfigFromEnv() (microservice.Config, error) {
	cfg := new(microservice.Config)
	err := env.Parse(cfg)
	return *cfg, err
}

package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq" // using postgres driver.
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"microservice"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdown(func() {
		cancel()
	})

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

	tracer, shutdown, err := initTracer("dummy_service", cfg.Traces.JaegerEndpoint)
	if err != nil {
		logger.Fatal("could not init tracer", zap.Error(err))
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		logger.Error("can't open connection to postgres", zap.Error(err))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("can't close connection to postgres", zap.Error(err))
		}
	}()
	if err := db.Ping(); err != nil {
		logger.Error("can't ping database", zap.Error(err))
		return
	}

	{ // configuring db pooling.
		db.SetMaxOpenConns(cfg.Database.MaxOpenConnections)
		db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
		db.SetConnMaxLifetime(cfg.Database.MaxConnLifetime)
	}

	app := microservice.New(logger, *cfg, tracer, db)

	logger.Info("servers shutdown err", zap.Error(app.Run(ctx)))
}

func initTracer(serviceName, jaegerURL string) (trace.Tracer, func(ctx context.Context) error, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	tracer := otel.Tracer(serviceName)

	return tracer, tp.Shutdown, nil
}

func gracefulShutdown(actions func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		actions()
	}()
}

func getConfigFromEnv() (*microservice.Config, error) {
	cfg := new(microservice.Config)
	err := env.Parse(cfg)
	return cfg, err
}

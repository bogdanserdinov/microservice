package microservice

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"microservice/database"
	"microservice/metrics"
	"microservice/pkg/http/server"
	"microservice/servers/private"
	"microservice/servers/public"
	"microservice/service"
)

type Config struct {
	Database struct {
		URL                string        `env:"URL,required"`
		MaxOpenConnections int           `env:"MAX_OPEN_CONNECTIONS" envDefault:"25"`
		MaxIdleConnections int           `env:"MAX_IDLE_CONNECTIONS" envDefault:"25"`
		MaxConnLifetime    time.Duration `env:"MAX_CONN_LIFETIME"    envDefault:"5m"`
	} `envPrefix:"DATABASE_"`

	Traces struct {
		JaegerEndpoint string `env:"JAEGER_ENDPOINT,required"`
	} `envPrefix:"TRACES_"`

	PublicServer  server.Config `envPrefix:"PUBLIC_SERVER_"`
	PrivateServer server.Config `envPrefix:"PRIVATE_SERVER_"`
}

type App struct {
	log *zap.Logger

	cfg Config

	dummy *service.Service

	publicServer  *public.Server
	privateServer *private.Server
}

func New(log *zap.Logger, cfg Config, tracer trace.Tracer, db *sql.DB) *App {
	app := &App{
		cfg: cfg,
		log: log,
	}

	registry := prometheus.NewRegistry()
	factory := promauto.With(registry)
	prom := metrics.New(factory)

	{ // service initialization.
		app.dummy = service.New(tracer, database.New(db))
	}

	{ // public server initialization.
		app.publicServer = public.New(
			cfg.PublicServer,
			log,
			prom,
			tracer,
			app.dummy,
		)
	}

	{ // private server initialization.
		app.privateServer = private.New(
			cfg.PrivateServer,
			log,
			registry,
			readinessProbe(db),
		)
	}

	return app
}

func (a *App) Run(ctx context.Context) error {
	var group errgroup.Group

	group.Go(func() error {
		err := a.publicServer.Run(ctx)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return errors.Wrap(err, "public server error")
	})
	group.Go(func() error {
		err := a.privateServer.Run(ctx)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return errors.Wrap(err, "private server error")
	})

	return group.Wait()
}

// readinessProbe is a helper function that returns readiness checker for given dependencies.
func readinessProbe(db *sql.DB) health.Checker {
	checker := health.NewChecker(
		health.WithCheck(health.Check{
			Name:    "database",
			Timeout: 2 * time.Second,
			Check:   db.PingContext,
		}),
	)

	return checker
}

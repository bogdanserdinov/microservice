package private

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"microservice/pkg/http/cors"
	"microservice/pkg/http/server"
)

type Server struct {
	cfg server.Config

	log *zap.Logger

	server *http.Server
}

func New(cfg server.Config, log *zap.Logger, registry *prometheus.Registry, checker health.Checker) *Server {
	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{})).Methods(http.MethodGet)
	router.Handle("/healthz", health.NewHandler(checker)).Methods(http.MethodGet)

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           cors.Allow(router),
		ReadHeaderTimeout: 2 * time.Second,
	}

	return &Server{
		cfg:    cfg,
		log:    log,
		server: server,
	}
}

func (server *Server) Run(ctx context.Context) error {
	server.log.Info("starting private API server")

	var group errgroup.Group

	group.Go(func() error {
		<-ctx.Done()
		return server.server.Shutdown(ctx)
	})
	group.Go(func() error {
		err := server.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return err
	})

	return group.Wait()
}

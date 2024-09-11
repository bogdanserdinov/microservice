package public

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"microservice/metrics"
	"microservice/pkg/http/cors"
	"microservice/pkg/http/server"
	"microservice/servers/public/controllers"
	"microservice/service"
)

type Server struct {
	cfg    server.Config
	log    *zap.Logger
	server *http.Server

	metrics *metrics.Metrics

	service *service.Service
}

func New(cfg server.Config, log *zap.Logger, service *service.Service, metrics *metrics.Metrics) *Server {
	server := &Server{
		cfg: cfg,
		server: &http.Server{
			Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
			ReadHeaderTimeout: 2 * time.Second,
		},
		log:     log,
		metrics: metrics,
		service: service,
	}

	server.initRoutes()

	return server
}

func (server *Server) initRoutes() {
	controller := controllers.NewDummy(server.log, server.service)

	router := mux.NewRouter()
	router.Use(server.ObserveHandlerDuration)
	apiRouter := router.PathPrefix("/v1").Subrouter()

	apiRouter.HandleFunc("/dummy", controller.Create).Methods(http.MethodPost)

	server.server.Handler = cors.Allow(router)
}

func (server *Server) ObserveHandlerDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			server.metrics.APIRequests.WithLabelValues().Observe(time.Since(start).Seconds())
		}()

		next.ServeHTTP(w, r.Clone(r.Context()))
	})
}

func (server *Server) Run(ctx context.Context) error {
	server.log.Info("starting public API server")

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

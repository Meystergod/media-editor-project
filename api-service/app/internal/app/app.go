package app

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/Meystergod/media-editor-project/api-service/internal/config"
	"github.com/Meystergod/media-editor-project/api-service/pkg/logging"
	"github.com/Meystergod/media-editor-project/api-service/pkg/metric"
	"github.com/Meystergod/media-editor-project/api-service/pkg/shutdown"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Info("router initializing")
	router := httprouter.New()

	logger.Info("heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	return App{
		cfg:    cfg,
		logger: logger,
		router: router,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, c := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(c)
	})

	return grp.Wait()
}

func (a *App) startHTTP(ctx context.Context) error {
	logger := a.logger.WithFields(map[string]interface{}{
		"IP":   a.cfg.HTTP.IP,
		"Port": a.cfg.HTTP.Port,
	})
	logger.Info("HTTP server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		a.logger.WithError(err).Fatal("failed to create listener")
	}

	logger.Info("cors initializing")
	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     a.cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   a.cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     a.cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: a.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     a.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              a.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go shutdown.Graceful(a.logger, []os.Signal{
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		os.Interrupt,
	}, a.httpServer)

	logger.Info("application initialized and started")

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warning("server shutdown")
		default:
			logger.Fatal(err)
		}
	}

	return err
}

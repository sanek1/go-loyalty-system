// Package app configures and runs application.
package app

import (
	"context"
	"go-loyalty-system/config"
	v1 "go-loyalty-system/internal/controller/v1"
	"go-loyalty-system/pkg/httpserver"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	//l := logger.New(cfg.Log.Level)
	ctx := context.Background()
	l, err := logging.NewZapLogger(zap.InfoLevel)
	if err != nil {
		panic(err)
	}

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.FatalCtx(ctx, "app - Run - postgres.New: %w", zap.Error(err))
	}
	defer pg.Close()

	// // Use case
	// translationUseCase := usecase.New(
	// 	repo.New(pg),
	// 	//webapi.New(),
	// )

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.InfoCtx(ctx, "app - Run - signal: "+s.String())
	case err = <-httpServer.Notify():
		l.InfoCtx(ctx, "app - Run - httpServer.Notify: %w", zap.Error(err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.ErrorCtx(ctx, "app - Run - httpServer.Shutdown: %w", zap.Error(err))
	}

}

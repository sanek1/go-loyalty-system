package app

import (
	"context"
	"go-loyalty-system/config"
	v1 "go-loyalty-system/internal/controller/http"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/internal/usecase/repo"
	"go-loyalty-system/pkg/httpserver"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	l, err := logging.NewZapLogger(zap.InfoLevel)
	if err != nil {
		panic(err)
	}

	// init db
	initPostgres()

	// Repository
	pg, err := postgres.NewPostgres(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.FatalCtx(ctx, "app - Run - postgres.New: %w", zap.Error(err))
	}
	defer pg.Close()

	// Use case
	gophermartUseCase := usecase.NewGopherMart(
		repo.NewUserRepo(pg, l),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, *gophermartUseCase, cfg)
	httpServer := httpserver.NewServer(handler, httpserver.Port(cfg.HTTP.Port))

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
	if err = httpServer.Shutdown(); err != nil {
		l.ErrorCtx(ctx, "app - Run - httpServer.Shutdown: %w", zap.Error(err))
	}
}

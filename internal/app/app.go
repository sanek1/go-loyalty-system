package app

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "go-loyalty-system/internal/controller/http"
	"go-loyalty-system/internal/controller/http/middleware"
	"go-loyalty-system/internal/controller/http/security"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/internal/usecase/repo"
	"go-loyalty-system/pkg/httpserver"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
)

type App struct {
	cfg        *config.Config
	logger     *logging.ZapLogger
	httpServer *http.Server
	postgres   *postgres.Postgres
}

// NewApp создает новый экземпляр приложения
func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()
	if cfg == nil {
		return nil, errors.New("config is required")
	}
	// Инициализация логгера
	log, err := logging.NewZapLogger(zap.InfoLevel)
	if err != nil {
		panic(err)
	}

	// Инициализация PostgreSQL
	initPostgres(cfg.PG.URL)

	pg, err := postgres.NewPostgres(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.FatalCtx(ctx, "app - Run - postgres.New: %w", zap.Error(err))
	}

	// Инициализация usecase
	uc := usecase.NewGopherMart(
		repo.NewUserRepo(pg, log, pg.Pool),
		log,
	)
	// middleware
	j := security.NewJwtToken(cfg.Jwt.EncryptionKey, *uc)
	a := middleware.NewAuthorizer(log)

	// init pool
	accrual := NewPoolController(*uc, cfg.Accrual.Accrual, log)
	startPool(accrual)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, *uc, cfg, j, accrual, a, log)
	httpServer := httpserver.NewServer(handler, httpserver.Port(cfg.HTTP.Port))

	return &App{
		cfg:        cfg,
		logger:     log,
		httpServer: httpServer.Server,
		postgres:   pg,
	}, nil
}

// Run запускает приложение
func (a *App) Run(ctx context.Context) error {
	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Канал для сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для ошибок сервера
	errChan := make(chan error, 1)

	// Ожидание сигнала или ошибки
	select {
	case <-ctx.Done():
		return a.shutdown()
	case sig := <-sigChan:
		a.logger.InfoCtx(ctx, "received signal", zap.String("signal", sig.String()))
		return a.shutdown()
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

// shutdown gracefully останавливает приложение
func (a *App) shutdown() error {
	a.logger.InfoCtx(context.Background(), "shutdown started")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Закрываем HTTP сервер
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	a.logger.InfoCtx(ctx, "shutdown completed")
	return nil
}

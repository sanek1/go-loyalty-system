package config

import (
	"context"
	"flag"
	"fmt"

	"os"
	"path/filepath"

	"go-loyalty-system/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type (
	// Config
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"logger"`
		PG      `yaml:"postgres"`
		Jwt     `yaml:"jwt"`
		Accrual `yaml:"accrual"`
	}

	App struct {
		Name    string ` yaml:"name"    env:"APP_NAME"`
		Version string ` yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Address string ` yaml:"Address" env:"RUN_ADDRESS"`
		Port    string ` yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string ` yaml:"log_level"   env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax int32  ` yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `                 env:"DATABASE_URI"`
	}

	Jwt struct {
		EncryptionKey string `json:"encryption_key" env:"AUTH_KEY"`
	}

	Accrual struct {
		Accrual string `json:"Accrual" env:"ACCRUAL_SYSTEM_ADDRESS"`
	}
)

func NewConfig() (*Config, error) {
	logger, _ := logging.NewZapLogger(0)
	//defer logger.Sync()
	cfg := &Config{}

	flag.StringVar(&cfg.HTTP.Address, "a", "", "RUN_ADDRESS")
	flag.StringVar(&cfg.PG.URL, "d", "", "Database URI")
	flag.StringVar(&cfg.Accrual.Accrual, "r", "", "Accrual address")
	flag.StringVar(&cfg.Jwt.EncryptionKey, "k", "", "Auth key")
	flag.Parse()

	executable, err := os.Executable()
	if err != nil {
		logger.ErrorCtx(context.Background(), "Failed to get executable path: %w", zap.Error(err))
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	configPath := filepath.Join(filepath.Dir(executable), "../config", "config.yaml")

	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		logger.ErrorCtx(context.Background(), "Failed to read config file: %w"+configPath, zap.Error(err))
		cfg.HTTP.Port = "8080"
		cfg.PG.PoolMax = 10
		cfg.Log.Level = "debug"
	}

	// Приоритет значений:
	// 1. Параметры командной строки
	// 2. Переменные окружения
	// 3. Значения из конфига
	// 4. Значения по умолчанию

	// Обработка переменных окружения
	if dbURI := os.Getenv("DATABASE_URI"); dbURI != "" {
		cfg.PG.URL = dbURI
	}

	if address := os.Getenv("RUN_ADDRESS"); address != "" {
		cfg.HTTP.Address = address
	}

	if key := os.Getenv("JWT_ENCRYPTION_KEY"); key != "" {
		cfg.Jwt.EncryptionKey = key
	}

	if accrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrual != "" {
		cfg.Accrual.Accrual = accrual
	}

	if cfg.HTTP.Address == "" {
		cfg.HTTP.Address = ":8080"
	}

	if cfg.PG.PoolMax == 0 {
		cfg.PG.PoolMax = 10
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "debug"
	}

	logger.InfoCtx(context.Background(), "Starting server with parameters",
		zap.String("address", cfg.HTTP.Address),
		zap.String("database", cfg.PG.URL),
		zap.String("accrual", cfg.Accrual.Accrual))

	return cfg, nil
}

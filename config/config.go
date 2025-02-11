package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
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
	cfg := &Config{}
	flag.StringVar(&cfg.HTTP.Address, "a", ":8081", "RUN_ADDRESS")
	flag.StringVar(&cfg.PG.URL, "d", "", "Database URI")
	flag.StringVar(&cfg.Accrual.Accrual, "r", ":8080", "Accrual address")
	flag.StringVar(&cfg.Jwt.EncryptionKey, "k", "", "Auth key")
	flag.Parse()

	if err := cleanenv.ReadConfig("config.yaml", cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	if key := os.Getenv("DATABASE_URI"); key != "" {
		cfg.URL = key
	}

	cfg.Jwt.EncryptionKey = os.Getenv("JWT_ENCRYPTION_KEY")

	if key := os.Getenv("RUN_ADDRESS"); key != "" {
		cfg.HTTP.Address = key
	}

	return cfg, nil
}

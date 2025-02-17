package app

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func initPostgres() {
	databaseURL, ok := os.LookupEnv("DATABASE_URI")
	if !ok || databaseURL == "" {
		log.Fatalf("migrate: environment variable not declared: DATABASE_URI")
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}
	migrationsPath := filepath.Join(currentDir, "migrations")
	log.Printf("Migrations path: %s", migrationsPath)

	for attempts > 0 {
		m, err = migrate.New("file:"+migrationsPath, databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d"+err.Error(), attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}
	defer m.Close()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}

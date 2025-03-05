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
	_countIterations = 15
)

func initPostgres(databaseURL string) {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	log.Printf("Migrations path: %s", currentDir)
	dir1 := "../../migrations"
	dir2 := "../migrations"

	for attempts > 0 {
		currentDir = filepath.Join(currentDir, dir1)
		m, err = migrate.New("file:"+currentDir, databaseURL)
		if err == nil {
			break
		}
		if attempts < _countIterations {
			currentDir = dir2
		}
		log.Printf("Migrations path: %s", currentDir)

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
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

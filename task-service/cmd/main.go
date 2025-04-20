package main

import (
	"fmt"
	"log/slog"
	"os"
	"task-service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)

	// TODO: изменить на os.Stderr
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Debug("Starting task-service...")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	migrations, err := migrate.New(
		"file://migrations",
		dbUrl,
	)

	if err != nil {
		log.Error("Failed to create migration instance", "error", err)
		return
	}

	if err := migrations.Up(); err != migrate.ErrNoChange {
		log.Error("Failed to apply migrations", "error", err)
		return
	}
	log.Info("Migrations applied successfully")

}

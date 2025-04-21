package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"task-service/internal/config"
	"task-service/internal/http/handlers/task/get"
	"task-service/internal/http/handlers/task/save"
	"task-service/internal/lib/logger/sl"
	"task-service/internal/lib/logger/sl/slogpretty"
	"task-service/internal/repo/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	// TODO: изменить на os.Stderr
	// log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log := setupPrettySlog()

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
		log.Error("Failed to create migration instance", sl.Error(err))
		os.Exit(1)
	}

	if err := migrations.Up(); err != migrate.ErrNoChange {
		log.Error("Failed to apply migrations", sl.Error(err))
		os.Exit(1)
	}
	log.Info("Migrations applied successfully")

	db, err := postgresql.NewDb(dbUrl)
	if err != nil {
		log.Error("Failed to connect to the database", sl.Error(err))
		os.Exit(1)
	}
	log.Info("Database connection established successfully")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/task", save.New(log, db))
	router.Get("/task/{id}", get.New(log, db))

	log.Info("Starting service", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IddleTimeout,
	}

	log.Info("HTTP server started", slog.String("address", cfg.HTTPServer.Address))
	log.Info("Waiting for requests...")
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start service")
	}
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stepan41k/Effective-Mobile/cmd/migrator"
	"github.com/stepan41k/Effective-Mobile/internal/app"
	"github.com/stepan41k/Effective-Mobile/internal/config"
	musicHandler "github.com/stepan41k/Effective-Mobile/internal/http-server/handlers/profile"
	musicService "github.com/stepan41k/Effective-Mobile/internal/service/profile"
	"github.com/stepan41k/Effective-Mobile/internal/storage/postgres"
	_ "github.com/stepan41k/Effective-Mobile/docs"
	"github.com/swaggo/http-swagger/v2"
)

// @title Effective Mobile Test API
// @version 0.1
// @description API Server for Effective Mobile application

// @host localhost:8082
// @BasePath /profile

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.DBName, os.Getenv("MY_DB_PASSWORD"), cfg.Storage.SSLMode)

	pool, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}
	service := musicService.New(pool, log)
	handler := musicHandler.New(service, log)

	storagePathForMigrator := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Storage.Username, os.Getenv("MY_DB_PASSWORD"), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.DBName, cfg.Storage.SSLMode)

	migrator.NewMigrator(storagePathForMigrator, os.Getenv("MY_MIGRATIONS_PATH"))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"), //The url pointing to API definition
	))

	router.Route("/profile", func(r chi.Router) {
		r.Post("/get", handler.GetProfiles(context.Background()))
		r.Delete("/delete", handler.DeleteProfile(context.Background()))
		r.Patch("/update", handler.UpdateProfile(context.Background()))
		r.Post("/create", handler.NewProfile(context.Background()))
	})

	log.Info("starting server")

	application := app.New(log, cfg, router)

	go func() {
		application.HTTPServer.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	log.Info("stopping application", slog.String("signal", signal.String()))

	application.HTTPServer.Stop(context.Background())

	postgres.Close(context.Background(), pool)

	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

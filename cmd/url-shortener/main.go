package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/svetlana-mel/url-shortener/internal/config"
	delete_url "github.com/svetlana-mel/url-shortener/internal/http-server/handlers/url/delete"
	get_url "github.com/svetlana-mel/url-shortener/internal/http-server/handlers/url/get"
	"github.com/svetlana-mel/url-shortener/internal/http-server/handlers/url/save"
	slog_lib "github.com/svetlana-mel/url-shortener/internal/lib/logger/slog"
	"github.com/svetlana-mel/url-shortener/internal/repository/sqlite"
)

const (
	ENV_LOCAL = "local"
	ENV_DEV   = "dev"
	ENV_PROD  = "prod"
)

func main() {
	cfg := config.NewConfig()

	log := setupLogger(ENV_LOCAL)

	log.Info("starting url-shortner", slog.String("env", cfg.Env))

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog_lib.AddErrorAtribute(err))
		os.Exit(1)
	}
	log.Info("storage init successfull")

	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Post("/", save.New(storage, log))
		r.Get("/{alias}", get_url.New(storage, log))
		r.Delete("/{alias}", delete_url.New(storage, log))
	})

	// router.Get("/alias", )
	// router.Delete("/alias", )

	// router.Get("/{alias}", redirect.New(storage, log))

	log.Info("starting server", slog.String("addr", cfg.Address))

	srv := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error run server", slog_lib.AddErrorAtribute(err))
	}

	log.Error("server stopped ")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case ENV_LOCAL:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_DEV:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_PROD:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

// TODO: config mechanism
// TODO: think about the deriving adapters

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/utilyre/contest/internal/adapters/handler"
	"github.com/utilyre/contest/internal/adapters/postgres"
	"github.com/utilyre/contest/internal/app/service"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	db, err := postgres.New(ctx)
	if err != nil {
		slog.Error("failed to establish database connection", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
			os.Exit(1)
		}
	}()

	accountRepo := postgres.NewAccountRepo(db)
	accountSvc := service.NewAccountService(accountRepo)

	r := chi.NewRouter()
	v1 := chi.NewRouter()

	v1.Route("/auth", func(r chi.Router) {
		handler := handler.NewAuthHandler(accountSvc)

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Get("/health", handler.NewHealthHandler().Check)
	r.Mount("/api/v1", v1)

	srv := &http.Server{
		Addr:    "localhost:3000",
		Handler: r,
	}
	slog.Info("starting http server", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start http server", "error", err)
		os.Exit(1)
	}
}

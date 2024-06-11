package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/utilyre/contest/internal/adapters/postgres"
	"github.com/utilyre/contest/internal/adapters/restful"
	"github.com/utilyre/contest/internal/app/service"
	"github.com/utilyre/contest/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("failed to load environment variables: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})))

	db, err := postgres.New(ctx, postgres.WithDSN(cfg.DBConnString))
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

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: restful.New(accountSvc),
	}
	slog.Info("starting http server", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start http server", "error", err)
		os.Exit(1)
	}
}

// TODO: config mechanism

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/utilyre/contest/internal/adapters/postgres"
	"github.com/utilyre/contest/internal/adapters/restful"
	"github.com/utilyre/contest/internal/app/service"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	db, err := postgres.New(ctx,
		postgres.WithUser("admin"),
		postgres.WithPassword("admin"),
		postgres.WithDBName("contest"),
		postgres.WithSSLMode("disable"),
	)
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
		Addr:    "localhost:3000",
		Handler: restful.New(accountSvc),
	}
	slog.Info("starting http server", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start http server", "error", err)
		os.Exit(1)
	}
}

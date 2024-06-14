package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"doc-app/config"
	"doc-app/database"
	"doc-app/jobs"

	"golang.org/x/exp/slog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config, err := config.NewConfig()
	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", err)
		return
	}

	database, err := database.NewDatabase(ctx, config)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to database", err)
		return
	}

	if err := database.Migrate(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to migrate database", err)
		return
	}

	client, err := jobs.Work(ctx, config, database)
	if err != nil {
		slog.ErrorContext(ctx, "failed to start worker", err)
		return
	}

	if err := client.Start(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to start worker", err)
		return
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Stop(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "failed to stop worker", err)
		return
	}
}

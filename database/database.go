//go:generate go run -tags gen ./sqlc

package database

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"doc-app/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	_ "doc-app/database/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*
var migrations embed.FS

type Database struct {
	config *pgxpool.Config
	*pgxpool.Pool
}

func NewDatabase(ctx context.Context, config config.Config) (Database, error) {
	cfg, err := pgxpool.ParseConfig(config.DATABASE_URL)
	if err != nil {
		return Database{}, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return Database{}, err
	}

	if err := pool.Ping(ctx); err != nil {
		return Database{}, err
	}

	return Database{
		config: cfg,
		Pool:   pool,
	}, nil
}

type gooseLogger struct{}

var _ goose.Logger = gooseLogger{}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	slog.Warn(fmt.Sprintf(format, v...))
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...))
}

func (d *Database) Migrate(ctx context.Context) error {
	db, err := goose.OpenDBWithDriver("pgx", d.config.ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetLogger(gooseLogger{})
	goose.SetBaseFS(migrations)
	return goose.UpContext(ctx, db, "migrations")
}

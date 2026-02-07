package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bekzat-kamen/startGO.git/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	connectionPath string
}

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Database)

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	slog.Info("connected to postgres")
	return db, nil
}

package repository

import (
	"log/slog"

	"github.com/bekzat-kamen/startGO.git/internal/config"
	"github.com/bekzat-kamen/startGO.git/internal/models"
)

type DB struct {
	connectionPath string
}

func (D DB) GetAll(courses []models.Course, err error) {
	//TODO implement me
	panic("implement me")
}

func NewPostgresDB(cfg *config.Config) (*DB, error) {
	db := &DB{}
	slog.Info("Connected to PostgreSQL")
	return db, nil
}

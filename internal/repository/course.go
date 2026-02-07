package repository

import (
	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type CourseRepo interface {
	GetAll() ([]models.Course, error)

	// TODO реализуй остальные методы
}

type PsgCourseRepo struct {
	db *sqlx.DB
}

func NewPsgCourseRepo(db *sqlx.DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (p *PsgCourseRepo) GetAll() ([]models.Course, error) {
	var courses []models.Course

	var query = `
	SELECT 	id, title, description, slug, price, duration,level, is_active, instructor_id, created_at, updated_at, deleted_at
	FROM courses
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC
	`
	err := p.db.Select(&courses, query)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

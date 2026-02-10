package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type CourseRepo interface {
	GetAll() ([]models.Course, error)
	GetByID(id int) (models.Course, error)
	DeleteByID(id int) error
	Create(input models.CreateCourse) (int, error)
}

type PsgCourseRepo struct {
	db *sqlx.DB
}

func NewPsgCourseRepo(db *sqlx.DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (p *PsgCourseRepo) Create(input models.CreateCourse) (int, error) {
	query := `
	INSERT INTO courses (
	title, description, slug, price, duration,level, is_active, instructor_id, created_at, updated_at
	) VALUES (
	          :title, :description, :slug, :price, :duration,:level, :is_active, :instructor_id, :created_at, :updated_at
	)
	RETURNING id
	`
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	stmt, err := p.db.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("Prepare query error: %w", err)
	}
	defer stmt.Close()
	var id int
	err = stmt.Get(&id, input)
	if err != nil {
		return 0, fmt.Errorf("Get query error: %w", err)
	}
	return id, nil
}

func (p *PsgCourseRepo) DeleteByID(id int) error {
	query := `UPDATE courses SET deleted_at = NOW(), updated_at=NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete course with id %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete course with id %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (p *PsgCourseRepo) GetByID(id int) (models.Course, error) {
	var course models.Course

	query := `SELECT id, title, description, slug, price, duration,level, is_active, instructor_id, created_at, updated_at, deleted_at 
			  FROM courses 
			  WHERE id=$1
			  AND deleted_at IS NULL
			  LIMIT 1`

	err := p.db.Get(&course, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return models.Course{}, models.ErrNotFound
		}

		return models.Course{}, fmt.Errorf("get course by id  %w", err)
	}
	return course, nil
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

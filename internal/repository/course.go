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
	Update(id int, input models.UpdateCourse) (int, error)
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
	var teacherExists bool
	checkTeacherQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND role = 'teacher')`
	if err := p.db.Get(&teacherExists, checkTeacherQuery, input.TeacherID); err != nil {
		return 0, fmt.Errorf("check teacher existence: %w", err)
	}
	if !teacherExists {
		return 0, models.ErrTeacherNotFound
	}

	var slugExists bool
	checkSlugQuery := `SELECT EXISTS (SELECT 1 FROM courses WHERE slug = $1 AND deleted_at IS NULL)`
	if err := p.db.Get(&slugExists, checkSlugQuery, input.Slug); err != nil {
		return 0, fmt.Errorf("check slug existence: %w", err)
	}
	if slugExists {
		return 0, models.ErrSlugAlreadyExists
	}

	query := `
	INSERT INTO courses (
	title, description, slug, price, duration, level, is_active, teacher_id, created_at, updated_at
	) VALUES (
	          :title, :description, :slug, :price, :duration, :level, :is_active, :teacher_id, :created_at, :updated_at
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
		return models.ErrCourseNotFound
	}
	return nil
}

func (p *PsgCourseRepo) GetByID(id int) (models.Course, error) {
	var course models.Course

	query := `SELECT id, title, description, slug, price, duration, level, is_active, teacher_id, created_at, updated_at, deleted_at
			  FROM courses 
			  WHERE id=$1
			  AND deleted_at IS NULL
			  LIMIT 1`

	err := p.db.Get(&course, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return models.Course{}, models.ErrCourseNotFound
		}

		return models.Course{}, fmt.Errorf("get course by id  %w", err)
	}
	return course, nil
}

func (p *PsgCourseRepo) GetAll() ([]models.Course, error) {
	var courses []models.Course

	var query = `
	SELECT id, title, description, slug, price, duration, level, is_active, teacher_id, created_at, updated_at, deleted_at
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

func (p *PsgCourseRepo) Update(id int, input models.UpdateCourse) (int, error) {
	current, err := p.GetByID(id)
	if err != nil {
		return 0, err
	}

	if input.Slug != nil && *input.Slug != current.Slug {
		var slugExists bool
		checkSlugQuery := `SELECT EXISTS (SELECT 1 FROM courses WHERE slug = $1 AND deleted_at IS NULL AND id <> $2)`
		if err := p.db.Get(&slugExists, checkSlugQuery, *input.Slug, id); err != nil {
			return 0, fmt.Errorf("check slug existence for update: %w", err)
		}
		if slugExists {
			return 0, models.ErrSlugAlreadyExists
		}
	}

	if input.TeacherID != nil {
		var teacherExists bool
		checkTeacherQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND role = 'teacher')`
		if err := p.db.Get(&teacherExists, checkTeacherQuery, *input.TeacherID); err != nil {
			return 0, fmt.Errorf("check teacher existence for update: %w", err)
		}
		if !teacherExists {
			return 0, models.ErrTeacherNotFound
		}
	}

	query := `
	UPDATE courses SET
	title = COALESCE($1, title),
	description = COALESCE($2, description),
	slug = COALESCE($3, slug),
	price = COALESCE($4, price),
	duration = COALESCE($5, duration),
	level = COALESCE($6, level),
	is_active = COALESCE($7, is_active),
	teacher_id = COALESCE($8, teacher_id),
	updated_at = NOW()
	WHERE id = $9 AND deleted_at IS NULL
	RETURNING id
	`

	var updatedID int
	err = p.db.QueryRow(
		query,
		input.Title,
		input.Description,
		input.Slug,
		input.Price,
		input.Duration,
		input.Level,
		input.IsActive,
		input.TeacherID,
		id,
	).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrCourseNotFound
		}
		return 0, fmt.Errorf("update course by id: %w", err)
	}

	return updatedID, nil
}

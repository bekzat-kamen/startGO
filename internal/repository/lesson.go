package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type LessonRepo interface {
	GetAll(ctx context.Context) ([]models.Lesson, error)
	GetByID(ctx context.Context, id int) (models.Lesson, error)
	DeleteByID(ctx context.Context, id int) error
	Create(ctx context.Context, input models.CreateLesson) (int, error)
	Update(ctx context.Context, id int, input models.UpdateLesson) (int, error)
}

type PsgLessonRepo struct {
	db *sqlx.DB
}

func NewPsgLessonRepo(db *sqlx.DB) *PsgLessonRepo {
	return &PsgLessonRepo{db: db}
}

func (p *PsgLessonRepo) GetAll(ctx context.Context) ([]models.Lesson, error) {
	var lessons []models.Lesson
	query := `
	SELECT id, course_id, title, content, video_url, duration, position, is_preview, created_at, updated_at, deleted_at
	FROM lessons
	WHERE deleted_at IS NULL
	ORDER BY position ASC, created_at ASC
	`
	if err := p.db.SelectContext(ctx, &lessons, query); err != nil {
		return nil, fmt.Errorf("get lessons: %w", err)
	}
	return lessons, nil
}

func (p *PsgLessonRepo) GetByID(ctx context.Context, id int) (models.Lesson, error) {
	var lesson models.Lesson
	query := `
	SELECT id, course_id, title, content, video_url, duration, position, is_preview, created_at, updated_at, deleted_at
	FROM lessons
	WHERE id = $1 AND deleted_at IS NULL
	LIMIT 1
	`
	if err := p.db.GetContext(ctx, &lesson, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Lesson{}, models.ErrLessonNotFound
		}
		return models.Lesson{}, fmt.Errorf("get lesson by id: %w", err)
	}
	return lesson, nil
}

func (p *PsgLessonRepo) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	var courseExists bool
	checkCourseQuery := `SELECT EXISTS (SELECT 1 FROM courses WHERE id = $1 AND deleted_at IS NULL)`
	if err := p.db.GetContext(ctx, &courseExists, checkCourseQuery, input.CourseID); err != nil {
		return 0, fmt.Errorf("check course existence: %w", err)
	}
	if !courseExists {
		return 0, models.ErrCourseNotFound
	}

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	query := `
	INSERT INTO lessons (course_id, title, content, video_url, duration, position, is_preview, created_at, updated_at)
	VALUES (:course_id, :title, :content, :video_url, :duration, :position, :is_preview, :created_at, :updated_at)
	RETURNING id
	`
	stmt, err := p.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare create lesson: %w", err)
	}
	defer stmt.Close()

	var id int
	if err := stmt.GetContext(ctx, &id, input); err != nil {
		return 0, fmt.Errorf("execute create lesson: %w", err)
	}
	return id, nil
}

func (p *PsgLessonRepo) Update(ctx context.Context, id int, input models.UpdateLesson) (int, error) {
	if _, err := p.GetByID(ctx, id); err != nil {
		return 0, err
	}

	if input.CourseID != nil {
		var courseExists bool
		checkCourseQuery := `SELECT EXISTS (SELECT 1 FROM courses WHERE id = $1 AND deleted_at IS NULL)`
		if err := p.db.GetContext(ctx, &courseExists, checkCourseQuery, *input.CourseID); err != nil {
			return 0, fmt.Errorf("check course existence for update: %w", err)
		}
		if !courseExists {
			return 0, models.ErrCourseNotFound
		}
	}

	query := `
	UPDATE lessons SET
	course_id = COALESCE($1, course_id),
	title = COALESCE($2, title),
	content = COALESCE($3, content),
	video_url = COALESCE($4, video_url),
	duration = COALESCE($5, duration),
	position = COALESCE($6, position),
	is_preview = COALESCE($7, is_preview),
	updated_at = NOW()
	WHERE id = $8 AND deleted_at IS NULL
	RETURNING id
	`
	var updatedID int
	err := p.db.QueryRowContext(
		ctx,
		query,
		input.CourseID,
		input.Title,
		input.Content,
		input.VideoURL,
		input.Duration,
		input.Position,
		input.IsPreview,
		id,
	).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrLessonNotFound
		}
		return 0, fmt.Errorf("update lesson by id: %w", err)
	}
	return updatedID, nil
}

func (p *PsgLessonRepo) DeleteByID(ctx context.Context, id int) error {
	query := `UPDATE lessons SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete lesson by id: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for delete lesson: %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrLessonNotFound
	}
	return nil
}

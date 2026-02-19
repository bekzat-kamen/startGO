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

type UserRepo interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	DeleteByID(ctx context.Context, id int) error
	Create(ctx context.Context, input models.CreateUser) (int, error)
	Update(ctx context.Context, id int, input models.UpdateUser) (int, error)
}

type PsgUserRepo struct {
	db *sqlx.DB
}

func NewPsgUserRepo(db *sqlx.DB) *PsgUserRepo {
	return &PsgUserRepo{db: db}
}

func (p *PsgUserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `
	SELECT id, full_name, email, password_hash, role, is_active, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	`
	if err := p.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}

func (p *PsgUserRepo) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	query := `
	SELECT id, full_name, email, password_hash, role, is_active, created_at, updated_at
	FROM users
	WHERE id = $1
	LIMIT 1
	`
	if err := p.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (p *PsgUserRepo) Create(ctx context.Context, input models.CreateUser) (int, error) {
	var emailExists bool
	checkEmailQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	if err := p.db.GetContext(ctx, &emailExists, checkEmailQuery, input.Email); err != nil {
		return 0, fmt.Errorf("check email existence: %w", err)
	}
	if emailExists {
		return 0, models.ErrEmailAlreadyExists
	}

	if input.Role == "" {
		input.Role = "student"
	}

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	query := `
	INSERT INTO users (full_name, email, password_hash, role, is_active, created_at, updated_at)
	VALUES (:full_name, :email, :password_hash, :role, :is_active, :created_at, :updated_at)
	RETURNING id
	`
	stmt, err := p.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare create user: %w", err)
	}
	defer stmt.Close()

	var id int
	if err := stmt.GetContext(ctx, &id, input); err != nil {
		return 0, fmt.Errorf("execute create user: %w", err)
	}
	return id, nil
}

func (p *PsgUserRepo) Update(ctx context.Context, id int, input models.UpdateUser) (int, error) {
	current, err := p.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	if input.Email != nil && *input.Email != current.Email {
		var emailExists bool
		checkEmailQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND id <> $2)`
		if err := p.db.GetContext(ctx, &emailExists, checkEmailQuery, *input.Email, id); err != nil {
			return 0, fmt.Errorf("check email existence for update: %w", err)
		}
		if emailExists {
			return 0, models.ErrEmailAlreadyExists
		}
	}

	query := `
	UPDATE users SET
	full_name = COALESCE($1, full_name),
	email = COALESCE($2, email),
	password_hash = COALESCE($3, password_hash),
	role = COALESCE($4, role),
	is_active = COALESCE($5, is_active),
	updated_at = NOW()
	WHERE id = $6
	RETURNING id
	`
	var updatedID int
	err = p.db.QueryRowContext(
		ctx,
		query,
		input.FullName,
		input.Email,
		input.PasswordHash,
		input.Role,
		input.IsActive,
		id,
	).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrUserNotFound
		}
		return 0, fmt.Errorf("update user by id: %w", err)
	}
	return updatedID, nil
}

func (p *PsgUserRepo) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user by id: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for delete user: %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, id int, input models.UpdateUser) (int, error)
	DeleteByID(ctx context.Context, id int) error
	Create(ctx context.Context, input models.CreateUser) (int, error)
}

var _ UserRepo = (*PsgUserRepo)(nil)

type PsgUserRepo struct {
	db *sqlx.DB
}

func NewPsgUserRepo(db *sqlx.DB) *PsgUserRepo {
	return &PsgUserRepo{db: db}
}

func (r *PsgUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := `
		SELECT
			id,
			full_name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r *PsgUserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT
			id,
			full_name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		ORDER BY created_at DESC
	`

	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}

	return users, nil
}

func (r *PsgUserRepo) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User

	query := `
		SELECT
			id,
			full_name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *PsgUserRepo) Create(ctx context.Context, input models.CreateUser) (int, error) {
	query := `
		INSERT INTO users (
			full_name,
			email,
			password_hash,
			created_at,
			updated_at
		)
		VALUES (
			:full_name,
			:email,
			:password_hash,
			:created_at,
			:updated_at
		)
		RETURNING id
	`

	input.CreatedAt = utils.Now()
	input.UpdatedAt = utils.Now()

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare create user: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.GetContext(ctx, &id, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				return 0, models.ErrUserAlreadyExists
			} // email exists
		}

		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (r *PsgUserRepo) Update(ctx context.Context, id int, input models.UpdateUser) (int, error) {
	if _, err := r.GetByID(ctx, id); err != nil {
		return 0, err
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
	err := r.db.QueryRowContext(
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, models.ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("update user: %w", err)
	}

	return updatedID, nil
}

func (r *PsgUserRepo) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

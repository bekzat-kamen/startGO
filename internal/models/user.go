package models

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	FullName     string    `db:"full_name" json:"full_name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	Role         string    `db:"role" json:"role"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateUser struct {
	FullName     string    `db:"full_name" json:"full_name" binding:"required"`
	Email        string    `db:"email" json:"email" binding:"required,email"`
	PasswordHash string    `db:"password_hash" json:"password_hash" binding:"required"`
	Role         string    `db:"role" json:"role"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UpdateUser struct {
	FullName     *string `db:"full_name" json:"full_name"`
	Email        *string `db:"email" json:"email"`
	PasswordHash *string `db:"password_hash" json:"password_hash"`
	Role         *string `db:"role" json:"role"`
	IsActive     *bool   `db:"is_active" json:"is_active"`
}

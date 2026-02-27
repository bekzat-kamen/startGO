package auth

import "github.com/bekzat-kamen/startGO.git/internal/models"

type TokenManager interface {
	NewAccessToken(user models.User) (string, int64, error)
}

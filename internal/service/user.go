package service

import (
	"context"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, input models.CreateUser) (int, error) {
	return s.repo.Create(ctx, input)
}

func (s *UserService) Update(ctx context.Context, id int, input models.UpdateUser) (int, error) {
	return s.repo.Update(ctx, id, input)
}

func (s *UserService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

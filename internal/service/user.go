package service

import (
	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetByID(id int) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Create(input models.CreateUser) (int, error) {
	return s.repo.Create(input)
}

func (s *UserService) Update(id int, input models.UpdateUser) (int, error) {
	return s.repo.Update(id, input)
}

func (s *UserService) DeleteByID(id int) error {
	return s.repo.DeleteByID(id)
}

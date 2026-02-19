package service

import (
	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/repository"
)

type LessonService struct {
	repo repository.LessonRepo
}

func NewLessonService(repo repository.LessonRepo) *LessonService {
	return &LessonService{repo: repo}
}

func (s *LessonService) GetAll() ([]models.Lesson, error) {
	return s.repo.GetAll()
}

func (s *LessonService) GetByID(id int) (models.Lesson, error) {
	return s.repo.GetByID(id)
}

func (s *LessonService) Create(input models.CreateLesson) (int, error) {
	return s.repo.Create(input)
}

func (s *LessonService) Update(id int, input models.UpdateLesson) (int, error) {
	return s.repo.Update(id, input)
}

func (s *LessonService) DeleteByID(id int) error {
	return s.repo.DeleteByID(id)
}

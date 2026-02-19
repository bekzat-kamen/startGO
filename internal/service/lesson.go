package service

import (
	"context"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/repository"
)

type LessonService struct {
	repo repository.LessonRepo
}

func NewLessonService(repo repository.LessonRepo) *LessonService {
	return &LessonService{repo: repo}
}

func (s *LessonService) GetAll(ctx context.Context) ([]models.Lesson, error) {
	return s.repo.GetAll(ctx)
}

func (s *LessonService) GetByID(ctx context.Context, id int) (models.Lesson, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *LessonService) Create(ctx context.Context, input models.CreateLesson) (int, error) {
	return s.repo.Create(ctx, input)
}

func (s *LessonService) Update(ctx context.Context, id int, input models.UpdateLesson) (int, error) {
	return s.repo.Update(ctx, id, input)
}

func (s *LessonService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}

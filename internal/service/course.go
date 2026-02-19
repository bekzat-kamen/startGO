package service

import (
	"context"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/repository"
)

type CourseService struct {
	repo repository.CourseRepo
}

func NewCourseService(repo repository.CourseRepo) *CourseService {
	return &CourseService{repo: repo}
}

func (cs *CourseService) Create(ctx context.Context, input models.CreateCourse) (int, error) {
	return cs.repo.Create(ctx, input)
}

func (cs *CourseService) GetAll(ctx context.Context) ([]models.Course, error) {
	return cs.repo.GetAll(ctx)
}

func (cs *CourseService) GetByID(ctx context.Context, id int) (models.Course, error) {
	return cs.repo.GetByID(ctx, id)
}

func (cs *CourseService) Update(ctx context.Context, id int, input models.UpdateCourse) (int, error) {
	return cs.repo.Update(ctx, id, input)
}

func (cs *CourseService) DeleteByID(ctx context.Context, id int) error {
	return cs.repo.DeleteByID(ctx, id)
}

package repository

import "github.com/bekzat-kamen/startGO.git/internal/models"

type CourseRepo interface {
	GetAll([]models.Course, error)

	// TODO реализуй остальные методы
}

type PsgCourseRepo struct {
	db *DB
}

func NewPsgCourseRepo(db *DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (p *PsgCourseRepo) GetAll() ([]models.Course, error) {
	//TODO исползовать db для получение данных и POSTGRES
	return []models.Course{
		{ID: 1, Name: "GO BASICS"},
		{ID: 2, Name: "SQL BASICS"},
		{ID: 3, Name: "JAVA BASICS"},
	}, nil
}

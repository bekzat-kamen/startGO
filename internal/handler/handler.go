package handler

import (
	"github.com/bekzat-kamen/startGO.git/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	courseService *service.CourseService
	lessonService *service.LessonService
	userService   *service.UserService
}

func NewHandler(cs *service.CourseService, ls *service.LessonService, us *service.UserService) *Handler {
	return &Handler{
		courseService: cs,
		lessonService: ls,
		userService:   us,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()

	r.GET("/courses", h.GetCourses)
	r.GET("/courses/:id", h.GetCourseByID) // localhost:8080/courses/@#@
	r.DELETE("/courses/:id", h.DeleteCourse)
	r.POST("/courses", h.CreateCourse)
	r.PUT("/courses/:id", h.UpdateCourse)

	r.GET("/lessons", h.GetLessons)
	r.GET("/lessons/:id", h.GetLessonByID)
	r.DELETE("/lessons/:id", h.DeleteLesson)
	r.POST("/lessons", h.CreateLesson)
	r.PUT("/lessons/:id", h.UpdateLesson)

	r.GET("/users", h.GetUsers)
	r.GET("/users/:id", h.GetUserByID)
	r.POST("/users", h.CreateUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)

	return r, nil
}

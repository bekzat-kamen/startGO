package handler

import (
	"net/http"

	"github.com/bekzat-kamen/startGO.git/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	courseService *service.CourseService

	// TODO добавить остальные сервисы
}

func NewHandler(cs *service.CourseService) *Handler {
	return &Handler{
		courseService: cs,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()
	r.GET("/courses", h.GetCourses)

	// other routes
	return r, nil
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.courseService.GetAll()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, courses)
	return
}

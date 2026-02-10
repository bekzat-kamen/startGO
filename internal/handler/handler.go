package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/bekzat-kamen/startGO.git/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	courseService *service.CourseService
}

func NewHandler(cs *service.CourseService) *Handler {
	return &Handler{
		courseService: cs,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, error) {
	r := gin.New()
	r.GET("/courses", h.GetCourses)
	r.GET("/courses/:id", h.GetCourseById)
	r.DELETE("/courses/:id", h.DeleteCourse)
	r.POST("/courses", h.CreateCourse)
	// other routes
	return r, nil
}

func (h *Handler) CreateCourse(c *gin.Context) {
	var input models.CreateCourse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.courseService.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})

}

func (h *Handler) DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.courseService.DeleteByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.courseService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, courses)
	return

}

func (h *Handler) GetCourseById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	course, err := h.courseService.GetCourseById(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": models.ErrNotFound.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

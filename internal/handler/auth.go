package handler

import (
	"errors"
	"net/http"

	"github.com/bekzat-kamen/startGO.git/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var input models.RegisterUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "failed to register user"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

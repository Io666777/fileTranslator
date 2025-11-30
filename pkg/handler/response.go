package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Исправь название - error это зарезервированное слово в Go
type errorResponse struct {
	Message string `json:"message"` // исправлено mesage -> message
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
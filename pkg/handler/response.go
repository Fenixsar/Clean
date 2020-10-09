package handler

import (
	"github.com/gin-gonic/gin"
	"gitlab.q123123.net/ligmar/boot"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	boot.Log.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

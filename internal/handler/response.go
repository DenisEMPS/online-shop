package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type CustomError struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, code int, message string) {
	log.Print(message)
	c.AbortWithStatusJSON(code, CustomError{Message: message})
}

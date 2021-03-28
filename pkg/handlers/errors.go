package handlers

import "github.com/gin-gonic/gin"

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) Error() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

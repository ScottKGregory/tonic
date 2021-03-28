package handlers

import "github.com/gin-gonic/gin"

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

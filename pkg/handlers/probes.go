package handlers

import "github.com/gin-gonic/gin"

type ProbeHandler struct {
}

func NewProbeHandler() *ProbeHandler {
	return &ProbeHandler{}
}

func (h *ProbeHandler) Health() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *ProbeHandler) Liveliness() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *ProbeHandler) Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

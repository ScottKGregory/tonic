package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/backends"
)

type ProbeHandler struct {
	backend backends.Backend
}

func NewProbeHandler(backend backends.Backend) *ProbeHandler {
	return &ProbeHandler{backend}
}

func (h *ProbeHandler) Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.backend.Ping(); err != nil {
			api.ErrorResponse(c, http.StatusInternalServerError, errors.New("Error connecting to backend"))
			return
		}

		api.SuccessResponse(c, "Healthy")
	}
}

func (h *ProbeHandler) Liveliness() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.backend.Ping(); err != nil {
			api.ErrorResponse(c, http.StatusInternalServerError, errors.New("Error connecting to backend"))
			return
		}

		api.SuccessResponse(c, "Alive")
	}
}

func (h *ProbeHandler) Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.backend.Ping(); err != nil {
			api.ErrorResponse(c, http.StatusInternalServerError, errors.New("Error connecting to backend"))
			return
		}

		api.SuccessResponse(c, "Ready")
	}
}

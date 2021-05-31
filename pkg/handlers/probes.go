package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/api"
	"github.com/scottkgregory/tonic/internal/backends"
)

type ProbeResponse struct {
	api.ResponseModel
	Data string
} //@Name ProbeResponse

type ProbeHandler struct {
	backend backends.Backend
}

func NewProbeHandler(backend backends.Backend) *ProbeHandler {
	return &ProbeHandler{backend}
}

// Health is the general health endpoint
// @Summary Get the health status of the service
// @Description Gets the health status of the service, returns error if database cannot be contacted
// @ID health
// @Tags probes
// @Produce json
// @Success 200 {object} ProbeResponse
// @Failure 400 {object} ProbeResponse
// @Failure 500 {object} ProbeResponse
// @Router /health [get]
func (h *ProbeHandler) Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.backend.Ping(); err != nil {
			api.ErrorResponse(c, http.StatusInternalServerError, errors.New("Error connecting to backend"))
			return
		}

		api.SuccessResponse(c, "Healthy")
	}
}

// Liveliness is the general liveliness endpoint
// @Summary Get the liveliness status of the service
// @Description Gets the liveliness status of the service, returns error if database cannot be contacted
// @ID liveliness
// @Tags probes
// @Produce json
// @Success 200 {object} ProbeResponse
// @Failure 400 {object} ProbeResponse
// @Failure 500 {object} ProbeResponse
// @Router /liveliness [get]
func (h *ProbeHandler) Liveliness() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.backend.Ping(); err != nil {
			api.ErrorResponse(c, http.StatusInternalServerError, errors.New("Error connecting to backend"))
			return
		}

		api.SuccessResponse(c, "Alive")
	}
}

// Readiness is the general readiness endpoint
// @Summary Get the readiness status of the service
// @Description Gets the readiness status of the service, returns error if database cannot be contacted
// @ID readiness
// @Tags probes
// @Produce json
// @Success 200 {object} ProbeResponse
// @Failure 400 {object} ProbeResponse
// @Failure 500 {object} ProbeResponse
// @Router /readiness [get]
func (h *ProbeHandler) Readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.backend.Ping(); err != nil {
			api.ErrorResponse(c, http.StatusInternalServerError, errors.New("Error connecting to backend"))
			return
		}

		api.SuccessResponse(c, "Ready")
	}
}

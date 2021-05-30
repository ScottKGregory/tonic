package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/api"
	"github.com/scottkgregory/tonic/internal/dependencies"
	"github.com/scottkgregory/tonic/internal/models"
	"github.com/scottkgregory/tonic/internal/services"
)

type PermissionsHandler struct {
	options *models.PermissionsOptions
}

func NewPermissionsHandler(options *models.PermissionsOptions) *PermissionsHandler {
	return &PermissionsHandler{options}
}

func (h *PermissionsHandler) ListPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewPermissionsService(log, h.options)

		perms, err := service.ListPermissions()
		api.SmartResponse(c, perms, err)
	}
}

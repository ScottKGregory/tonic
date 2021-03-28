package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/services"
)

type PermissionsHandler struct {
}

func NewPermissionsHandler() *PermissionsHandler {
	return &PermissionsHandler{}
}

func (h *PermissionsHandler) ListPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewPermissionsService(log)

		perms, err := service.ListPermissions()
		api.SmartResponse(c, perms, err)
	}
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

type ListPermissionsResponse struct {
	api.ResponseModel
	Data []string
} //@Name ListPermissionsResponse

type PermissionsHandler struct {
	options *models.PermissionsOptions
}

func NewPermissionsHandler(options *models.PermissionsOptions) *PermissionsHandler {
	return &PermissionsHandler{options}
}

// ListPermissions lists all permissions using the configured backend
// @Summary List all registered permissions
// @Description Lists all registered permissions
// @ID list-permissions
// @Tags permissions
// @Accept json
// @Produce json
// @Success 200 {object} ListPermissionsResponse
// @Failure 400 {object} ListPermissionsResponse
// @Failure 500 {object} ListPermissionsResponse
// @Router /api/permissions [get]
func (h *PermissionsHandler) ListPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewPermissionsService(log, h.options)

		perms, err := service.ListPermissions()
		api.SmartResponse(c, perms, err)
	}
}

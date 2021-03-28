package handlers

import "github.com/gin-gonic/gin"

type PermissionsHandler struct {
}

func NewPermissionsHandler() *PermissionsHandler {
	return &PermissionsHandler{}
}

func (h *PermissionsHandler) ListPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

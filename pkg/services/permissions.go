package services

import (
	"strings"

	"github.com/scottkgregory/tonic/pkg/models"
)

var TonicPermissions = []string{
	"GetToken",
	"ManageUsers",
}

type PermissionService struct {
	defaults    []string
	permissions []string
}

func NewPermissionService(options *models.Permissions) *PermissionService {
	// log := helpers.GetLogger()

	var perms []string
	for _, p := range append(TonicPermissions, options.List...) {
		perms = append(perms, strings.ToLower(p))
	}

	var defaults []string
	for _, p := range options.Defaults {
		defaults = append(defaults, strings.ToLower(p))
	}

	return &PermissionService{defaults: defaults, permissions: perms}
}

func (s *PermissionService) DefaultPermissions() []string {
	return s.defaults
}

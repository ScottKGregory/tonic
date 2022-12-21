package services

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/pkg/models"
)

type PermissionsService struct {
	log         *zerolog.Logger
	permissions []string
	config      *models.PermissionsConfig
}

// NewPermissionService initialises a new PermissionService based on the config supplied
func NewPermissionsService(log *zerolog.Logger, config *models.PermissionsConfig) *PermissionsService {
	return &PermissionsService{
		log: log,
		permissions: append([]string{
			"users:create:*",
			"users:update:*",
			"users:delete:*",
			"users:get:*",
			"users:list:*",
			"token:get:*",
			"permissions:list:*",
		}, config.Custom...),
		config: config,
	}
}

func (s *PermissionsService) ListPermissions() (out []string, err error) {
	for _, perm := range s.permissions {
		out = append(out, strings.ToLower(perm))
	}

	return out, nil
}

func (s *PermissionsService) DefaultPermissions() (out []string) {
	for _, perm := range s.config.Default {
		out = append(out, strings.ToLower(perm))
	}

	return out
}

func ValidatePermissions(perms ...string) (valid bool, messages map[string]string) {
	valid = true
	messages = map[string]string{}
	for _, p := range perms {
		if len(strings.Split(p, ":")) != 3 {
			messages[p] = "must have 3 parts"
			valid = false
		}
	}

	return valid, messages
}

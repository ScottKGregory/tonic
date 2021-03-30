package services

import (
	"strings"

	"github.com/rs/zerolog"
)

type PermissionsService struct {
	log         *zerolog.Logger
	permissions []string
}

// NewPermissionService initialises a new PermissionService based on the options supplied
func NewPermissionsService(log *zerolog.Logger) *PermissionsService {
	return &PermissionsService{log: log,
		permissions: []string{
			"users:create:*",
			"users:update:*",
			"users:delete:*",
			"users:get:*",
			"users:list:*",
			"token:get:*",
			"permissions:list:*",
		}}
}

func (s *PermissionsService) ListPermissions() (out []string, err error) {
	for _, perm := range s.permissions {
		out = append(out, strings.ToLower(perm))
	}

	return out, nil
}

func (s *PermissionsService) DefaultPermissions() (out []string) {
	for _, perm := range s.permissions {
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

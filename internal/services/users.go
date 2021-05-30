package services

import (
	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/internal/api/errors"
	"github.com/scottkgregory/tonic/internal/backends"
	"github.com/scottkgregory/tonic/internal/constants"
	"github.com/scottkgregory/tonic/internal/helpers"
	"github.com/scottkgregory/tonic/pkg/models"
)

type UserService struct {
	log     *zerolog.Logger
	backend backends.Backend
}

// NewUserService initialises a new UserService based on the options supplied
func NewUserService(log *zerolog.Logger, backend backends.Backend) *UserService {
	return &UserService{log, backend}
}

// CreateUser uses the configured backend to create the supplied user after having validted it
func (s *UserService) CreateUser(in *models.User) (out *models.User, err error) {
	valid, messages := s.isValidUser(in)
	if !valid {
		return out, errors.NewValidationError(messages)
	}

	return s.backend.CreateUser(in)
}

// CreateUser uses the configured backend to update the supplied user after having validted it
func (s *UserService) UpdateUser(in *models.User, sub string) (out *models.User, err error) {
	valid, messages := s.isValidUser(in)
	if !valid {
		return out, errors.NewValidationError(messages)
	}

	if in.Claims.Subject != sub {
		messages["claims.subject"] = "Field does not match supplied param"
		return out, errors.NewValidationError(messages)
	}

	out, err = s.backend.UpdateUser(in)
	if err != nil {
		return out, err
	}

	if out == nil {
		messages[constants.GlobalKey] = "User does not exist to update"
		return out, errors.NewValidationError(messages)
	}

	return out, err
}

// DeleteUser uses the configured backend to mark the user as deleted
func (s *UserService) DeleteUser(sub string) error {
	user, err := s.GetUser(sub)
	if err != nil {
		return err
	}

	user.Deleted = true

	_, err = s.UpdateUser(user, sub)
	return err
}

// GetUser uses the configured backend to get a single user based on it's subject claim
func (s *UserService) GetUser(sub string) (out *models.User, err error) {
	out, err = s.backend.GetUser(sub)
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, errors.NewNotFoundError(sub)
	}

	return out, nil
}

// ListUsers uses the configured backend to list all users
func (s *UserService) ListUsers() (out []*models.User, err error) {
	return s.backend.ListUsers()
}

func (s *UserService) isValidUser(user *models.User) (valid bool, messages map[string]string) {
	valid = true
	messages = make(map[string]string)
	if helpers.IsEmptyOrWhitespace(user.Claims.Subject) {
		valid = false
		messages["claims.subject"] = "This field is missing"
	}

	return valid, messages
}

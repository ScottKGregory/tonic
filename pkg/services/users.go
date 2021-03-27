package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/helpers"
	"github.com/scottkgregory/tonic/pkg/models"
)

type UserService struct {
	backend     *backends.Mongo
	permService *PermissionService
}

func NewUserService(backendOptions *models.Backend, permService *PermissionService) *UserService {
	// log := helpers.GetLogger()

	return &UserService{
		backend:     backends.NewMongoBackend(backendOptions),
		permService: permService,
	}
}

func (s *UserService) UpdateUser(c *gin.Context) {
	log := helpers.GetLogger()
	id := c.Param(constants.IDParam)
	if helpers.IsEmptyOrWhitespace(id) {
		log.Warn().Msg("No ID supplied in url")
		helpers.APIErrorResponse(c, http.StatusBadRequest, "No ID supplied in url")
		return
	}

	var user *models.User
	err := c.BindJSON(&user)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		helpers.APIErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// msgs := make(map[string]string)
	// Validation

	err = s.backend.SaveUser(log, user)
	if err != nil {
		log.Error().Err(err).Msg("Error saving user")
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error saving user")
		return
	}

	helpers.APISuccessResponse(c, user)
}

func (s *UserService) GetUser(c *gin.Context) {
	log := helpers.GetLogger()
	id := c.Param(constants.IDParam)
	if helpers.IsEmptyOrWhitespace(id) {
		log.Warn().Msg("No ID supplied in url")
		helpers.APIErrorResponse(c, http.StatusBadRequest, "No ID supplied in url")
		return
	}

	user, err := s.backend.GetUserByID(log, id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user")
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error getting user")
		return
	}

	helpers.APISuccessResponse(c, user)
}

func (s *UserService) ListUsers(c *gin.Context) {
	log := helpers.GetLogger()

	users, err := s.backend.ListUsers(log)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user")
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error getting user")
	}

	helpers.APISuccessResponse(c, users)
}

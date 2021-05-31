package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/api"
	"github.com/scottkgregory/tonic/internal/backends"
	"github.com/scottkgregory/tonic/internal/constants"
	"github.com/scottkgregory/tonic/internal/dependencies"
	"github.com/scottkgregory/tonic/internal/services"
	"github.com/scottkgregory/tonic/pkg/models"
)

type UserResponse struct {
	api.ResponseModel
	Data models.User
} //@Name UserResponse

type ListUserResponse struct {
	api.ResponseModel
	Data models.User
} //@Name ListUserResponse

type UserHandler struct {
	backend backends.Backend
}

func NewUserHandler(backend backends.Backend) *UserHandler {
	return &UserHandler{backend}
}

// CreateUser creates a user using the configured backend
// @Summary Create a single user
// @Description Creates a single user
// @ID create-user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Failure 500 {object} UserResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		model := &models.User{}
		err := c.Bind(model)
		if err != nil {
			log.Error().Err(err).Msg("Error binding model")
			api.ValidationErrorResponse(c)
			return
		}

		out, err := service.CreateUser(model)
		api.SmartResponse(c, out, err)
	}
}

// UpdateUser updates a user using the configured backend
// @Summary Update a single user
// @Description Updates the supplied user
// @ID update-user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Failure 500 {object} UserResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		model := &models.User{}
		err := c.Bind(model)
		if err != nil {
			log.Error().Err(err).Msg("Error binding model")
			api.ValidationErrorResponse(c)
			return
		}

		out, err := service.UpdateUser(model, c.Param(constants.IDParam))
		api.SmartResponse(c, out, err)
	}
}

// DeleteUser deletes a user using the configured backend
// @Summary Delete a single user
// @Description Deletes a single user
// @ID delete-user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Failure 500 {object} UserResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		err := service.DeleteUser(c.Param(constants.IDParam))
		api.SmartResponse(c, nil, err)
	}
}

// GetUser gets a single user using the configured backend
// @Summary Get a single user
// @Description Gets a user by ID
// @ID get-user-by-id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Failure 500 {object} UserResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		out, err := service.GetUser(c.Param(constants.IDParam))
		api.SmartResponse(c, out, err)
	}
}

// ListUsers lists all users using the configured backend
// @Summary List all registered users
// @Description Lists all registered users
// @ID list-users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} ListUserResponse
// @Failure 400 {object} ListUserResponse
// @Failure 500 {object} ListUserResponse
// @Router /api/users [get]
func (h *UserHandler) ListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		out, err := service.ListUsers()
		api.SmartResponse(c, out, err)
	}
}

// Me returns the currently authed user
// @Summary get the currently authed user
// @Description Gets the details of the currently authed user
// @ID me
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 400 {object} UserResponse
// @Failure 500 {object} UserResponse
// @Router /me [get]
func (h *UserHandler) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := dependencies.GetUser(c)
		if !ok {
			api.UnauthorisedResponse(c)
			return
		}

		c.Get(constants.UserKey)
		api.SmartResponse(c, user, nil)
	}
}

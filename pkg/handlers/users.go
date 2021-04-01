package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

type UserHandler struct {
	backend backends.Backend
}

func NewUserHandler(backend backends.Backend) *UserHandler {
	return &UserHandler{backend}
}

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

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		err := service.DeleteUser(c.Param(constants.IDParam))
		api.SmartResponse(c, nil, err)
	}
}

func (h *UserHandler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		out, err := service.GetUser(c.Param(constants.IDParam))
		api.SmartResponse(c, out, err)
	}
}

func (h *UserHandler) ListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		service := services.NewUserService(log, h.backend)

		out, err := service.ListUsers()
		api.SmartResponse(c, out, err)
	}
}

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

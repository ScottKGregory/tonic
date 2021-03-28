package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

const (
	errorRedirect  = "/error/500"
	logoutRedirect = "/"
	loginRedirect  = "/"
)

type AuthHandler struct {
	backend backends.Backend
	options *models.Auth
}

func NewAuthHandler(backend backends.Backend, options *models.Auth) *AuthHandler {
	return &AuthHandler{backend, options}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, h.backend)
		authService := services.NewAuthService(log, userService, h.options)

		url, err := authService.Login("")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, errorRedirect)
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, h.backend)
		authService := services.NewAuthService(log, userService, h.options)

		token, err := authService.Callback(
			c.Request.Context(),
			"",
			c.Query("state"),
			c.Query("code"),
			c.Query("error"),
			c.Query("error_description"),
		)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, errorRedirect)
		}

		c.SetCookie(
			h.options.Cookie.Name,
			string(token),
			int(h.options.JWT.Duration)*60,
			h.options.Cookie.Path,
			h.options.Cookie.Domain,
			h.options.Cookie.Secure,
			h.options.Cookie.HttpOnly,
		)

		c.Redirect(http.StatusTemporaryRedirect, loginRedirect) // Use a return URL in state
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(
			h.options.Cookie.Name,
			"",
			-1,
			h.options.Cookie.Path,
			h.options.Cookie.Domain,
			h.options.Cookie.Secure,
			h.options.Cookie.HttpOnly,
		)
		c.Redirect(http.StatusTemporaryRedirect, logoutRedirect)
	}
}

func (h *AuthHandler) Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, h.backend)
		authService := services.NewAuthService(log, userService, h.options)

		token, err := authService.Token(c.GetString(constants.SubjectKey))
		api.SmartResponse(c, token, err)
	}
}
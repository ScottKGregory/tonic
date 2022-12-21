package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/api/errors"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

const (
	errorRedirect    = "/error/500"
	unauthedRedirect = "/error/401"
	logoutRedirect   = "/"
	loginRedirect    = "/"
)

type AuthHandler struct {
	backend    backends.Backend
	config     *models.AuthConfig
	permConfig *models.PermissionsConfig
}

func NewAuthHandler(backend backends.Backend, config *models.AuthConfig, permConfig *models.PermissionsConfig) *AuthHandler {
	return &AuthHandler{backend, config, permConfig}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, h.backend)
		permService := services.NewPermissionsService(log, h.permConfig)
		authService := services.NewAuthService(log, userService, permService, h.config)

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
		permService := services.NewPermissionsService(log, h.permConfig)
		authService := services.NewAuthService(log, userService, permService, h.config)

		token, err := authService.Callback(c,
			c.Request.Context(),
			"",
			c.Query("state"),
			c.Query("code"),
			c.Query("error"),
			c.Query("error_description"),
		)
		if errors.Is(err, &errors.UnauthorisedErr{}) {
			c.Redirect(http.StatusTemporaryRedirect, unauthedRedirect)
			return
		} else if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, errorRedirect)
			return
		}

		c.SetCookie(
			h.config.Cookie.Name,
			string(token),
			int(h.config.JWT.Duration)*60,
			h.config.Cookie.Path,
			h.config.Cookie.Domain,
			h.config.Cookie.Secure,
			h.config.Cookie.HttpOnly,
		)

		c.Redirect(http.StatusTemporaryRedirect, loginRedirect) // Use a return URL in state
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(
			h.config.Cookie.Name,
			"",
			-1,
			h.config.Cookie.Path,
			h.config.Cookie.Domain,
			h.config.Cookie.Secure,
			h.config.Cookie.HttpOnly,
		)
		c.Redirect(http.StatusTemporaryRedirect, logoutRedirect)
	}
}

func (h *AuthHandler) Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, h.backend)
		permService := services.NewPermissionsService(log, h.permConfig)
		authService := services.NewAuthService(log, userService, permService, h.config)

		token, err := authService.Token(c.Request.Context(), c.GetString(constants.SubjectKey))
		api.SmartResponse(c, token, err)
	}
}

package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

const (
	bearerPrefix = "bearer "
)

func Authed(backend backends.Backend,
	cookieConfig *models.CookieConfig,
	jwtConfig *models.JWTConfig,
	authConfig *models.AuthConfig,
	permissionConfig *models.PermissionsConfig,
	cancel bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, backend)
		permService := services.NewPermissionsService(log, permissionConfig)
		authService := services.NewAuthService(log, userService, permService, authConfig)

		header := c.GetHeader(constants.Authorization)
		token, err := c.Cookie(cookieConfig.Name)
		if err != nil {
			if header == "" {
				retErr(c, cookieConfig, cancel)
				return
			}

			if !strings.HasPrefix(strings.ToLower(header), bearerPrefix) {
				retErr(c, cookieConfig, cancel)
				return
			}

			token = header[len(bearerPrefix)-1:]
			c.Set(constants.AuthMethodKey, constants.Bearer)
		} else {
			c.Set(constants.AuthMethodKey, constants.Cookie)
		}

		valid, validToken := authService.Verify(token)
		if !valid {
			retErr(c, cookieConfig, cancel)
			return
		}

		subject := validToken.Subject()
		expiry := validToken.Expiration()

		l := log.With().Str("user", subject).Logger()
		c.Set(constants.LoggerKey, &l)
		log = &l

		if time.Until(expiry) <= (time.Duration(jwtConfig.Duration)*time.Minute)/2 {
			l.Debug().Msg("Renewing auth")
			newToken, err := authService.Token(c.Request.Context(), subject)
			if err != nil {
				retErr(c, cookieConfig, cancel)
				return
			}

			c.SetCookie(
				cookieConfig.Name,
				newToken.Token,
				int(jwtConfig.Duration)*int(time.Minute),
				cookieConfig.Path,
				cookieConfig.Domain,
				cookieConfig.Secure,
				cookieConfig.HttpOnly,
			)
		}

		user, err := userService.GetUser(c.Request.Context(), subject)
		if err != nil {
			retErr(c, cookieConfig, cancel)
			return
		}

		c.Set(constants.Authed, true)
		c.Set(constants.SubjectKey, subject)
		c.Set(constants.UserKey, user)

		c.Next()
	}
}

func retErr(c *gin.Context, cookieConfig *models.CookieConfig, cancel bool) {
	if cancel {
		c.SetCookie(cookieConfig.Name, "", -1, cookieConfig.Path, cookieConfig.Domain, cookieConfig.Secure, cookieConfig.HttpOnly)
		api.UnauthorisedResponse(c)
		c.Abort()
	}

	c.Set(constants.Authed, false)
	c.Next()
}

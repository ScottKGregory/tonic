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
	cookieOptions *models.CookieOptions,
	jwtOptions *models.JWTOptions,
	authOptions *models.AuthOptions,
	permissionOptions *models.PermissionsOptions,
	cancel bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, backend)
		permService := services.NewPermissionsService(log, permissionOptions)
		authService := services.NewAuthService(log, userService, permService, authOptions)

		header := c.GetHeader(constants.Authorization)
		token, err := c.Cookie(cookieOptions.Name)
		if err != nil {
			if header == "" {
				retErr(c, cookieOptions, cancel)
				return
			}

			if !strings.HasPrefix(strings.ToLower(header), bearerPrefix) {
				retErr(c, cookieOptions, cancel)
				return
			}

			token = header[len(bearerPrefix)-1:]
			c.Set(constants.AuthMethodKey, constants.Bearer)
		} else {
			c.Set(constants.AuthMethodKey, constants.Cookie)
		}

		valid, validToken := authService.Verify(token)
		if !valid {
			retErr(c, cookieOptions, cancel)
			return
		}

		subject := validToken.Subject()
		expiry := validToken.Expiration()

		l := log.With().Str("user", subject).Logger()
		c.Set(constants.LoggerKey, &l)
		log = &l

		if time.Until(expiry) <= (time.Duration(jwtOptions.Duration)*time.Minute)/2 {
			l.Debug().Msg("Renewing auth")
			newToken, err := authService.Token(c.Request.Context(), subject)
			if err != nil {
				retErr(c, cookieOptions, cancel)
				return
			}

			c.SetCookie(
				cookieOptions.Name,
				newToken.Token,
				int(jwtOptions.Duration)*int(time.Minute),
				cookieOptions.Path,
				cookieOptions.Domain,
				cookieOptions.Secure,
				cookieOptions.HttpOnly,
			)
		}

		user, err := userService.GetUser(c.Request.Context(), subject)
		if err != nil {
			retErr(c, cookieOptions, cancel)
			return
		}

		c.Set(constants.Authed, true)
		c.Set(constants.SubjectKey, subject)
		c.Set(constants.UserKey, user)

		c.Next()
	}
}

func retErr(c *gin.Context, cookieOptions *models.CookieOptions, cancel bool) {
	if cancel {
		c.SetCookie(cookieOptions.Name, "", -1, cookieOptions.Path, cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly)
		api.UnauthorisedResponse(c)
		c.Abort()
	}

	c.Set(constants.Authed, false)
	c.Next()
}

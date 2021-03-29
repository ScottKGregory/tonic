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

func Authed(backend backends.Backend, cookieOptions *models.Cookie, jwtOptions *models.JWT, authOptions *models.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, backend)
		authService := services.NewAuthService(log, userService, authOptions)

		header := c.GetHeader(constants.Authorization)
		token, err := c.Cookie(cookieOptions.Name)
		if err != nil {
			if header == "" {
				retErr(c, cookieOptions)
				return
			}

			if !strings.HasPrefix(header, bearerPrefix) {
				retErr(c, cookieOptions)
				return
			}

			token = strings.TrimPrefix(header, bearerPrefix)
			c.Set(constants.AuthMethodKey, constants.Bearer)
		} else {
			c.Set(constants.AuthMethodKey, constants.Cookie)
		}

		valid, validToken := authService.Verify(token)
		if !valid {
			retErr(c, cookieOptions)
			return
		}

		subject := validToken.Subject()
		expiry := validToken.Expiration()

		l := log.With().Str("user", subject).Logger()
		c.Set(constants.LoggerKey, &l)
		log = &l

		if time.Until(expiry) <= (time.Duration(jwtOptions.Duration)*time.Minute)/2 {
			log.Debug().Msg("Renewing auth")
			newToken, err := authService.Token(subject)
			if err != nil {
				retErr(c, cookieOptions)
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

		c.Set(constants.Authed, true)
		c.Set(constants.SubjectKey, subject)

		c.Next()
	}
}

func retErr(c *gin.Context, cookieOptions *models.Cookie) {
	c.SetCookie(cookieOptions.Name, "", -1, cookieOptions.Path, cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly)
	api.UnauthorisedResponse(c)
	c.Abort()
}

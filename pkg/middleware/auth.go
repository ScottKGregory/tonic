package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/helpers"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/scottkgregory/tonic/pkg/services"
)

func Authed(cookieOptions *models.Cookie, jwtOptions *models.JWT, service *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := helpers.GetLogger(c)
		header := c.GetHeader(constants.Authorization)
		token, err := c.Cookie(cookieOptions.Name)
		if err != nil {

			if header == "" {
				log.Warn().Msg("Cookie or authorization header not present")
				c.SetCookie(cookieOptions.Name, "", -1, cookieOptions.Path, cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly)
				helpers.APIErrorResponse(c, http.StatusUnauthorized, "Unauthorized, authorization header or cookie required")
				c.Abort()
				return
			}

			if !strings.HasPrefix(header, services.BearerPrefix) {
				helpers.APIErrorResponse(c, http.StatusUnauthorized, fmt.Sprintf("Bearer token header must begin with '%s'", services.BearerPrefix))
				c.Abort()
				return
			}

			token = strings.TrimPrefix(header, services.BearerPrefix)
			c.Set(constants.AuthMethodKey, constants.Bearer)
		} else {
			c.Set(constants.AuthMethodKey, constants.Cookie)
		}

		valid, subject, expiry, perms := service.Verify(token, log)
		if !valid {
			log.Warn().Msg("Auth rejected")
			c.SetCookie(cookieOptions.Name, "", -1, cookieOptions.Path, cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly)
			helpers.APIErrorResponse(c, http.StatusUnauthorized, "Unauthorized, authorization header or cookie required")
			c.Abort()
			return
		}

		l := log.With().Str("user", subject).Logger()
		c.Set(constants.LoggerKey, &l)
		log = &l

		if time.Until(expiry) <= (time.Duration(jwtOptions.Duration)*time.Minute)/2 {
			log.Debug().Msg("Renewing auth")
			token, err = service.Renew(token, log)
			if err != nil {
				log.Error().Err(err).Msg("Error renewing token")
				c.SetCookie(cookieOptions.Name, "", -1, cookieOptions.Path, cookieOptions.Domain, cookieOptions.Secure, cookieOptions.HttpOnly)
				helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error renewing auth")
				c.Abort()
				return
			}

			c.SetCookie(
				cookieOptions.Name,
				string(token),
				int(jwtOptions.Duration)*int(time.Minute),
				cookieOptions.Path,
				cookieOptions.Domain,
				cookieOptions.Secure,
				cookieOptions.HttpOnly,
			)
		}

		c.Set(constants.SubjectKey, subject)
		c.Set(constants.PermissionsKey, perms)

		c.Next()
	}
}

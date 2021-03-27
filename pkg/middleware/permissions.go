package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/helpers"
)

func Any(handler gin.HandlerFunc, required ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := helpers.GetLogger(c)

		userPerms := c.GetStringSlice(constants.PermissionsKey)
		for _, perm := range userPerms {
			if contains(required, perm) {
				log.Info().Msg("Permissions check passed")
				handler(c)
				return
			}
		}

		log.Warn().Strs("required", required).Msg("User does not have any of the required permissions")
		helpers.APIErrorResponse(c, http.StatusForbidden, "User does not have any of the required permissions")
		c.Abort()
	}
}

func All(handler gin.HandlerFunc, required ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := helpers.GetLogger(c)

		matched := 0
		userPerms := c.GetStringSlice(constants.PermissionsKey)
		for _, perm := range userPerms {
			if contains(required, perm) {
				matched += 1
			}
		}

		if matched == len(required) {
			log.Info().Msg("Permissions check passed")
			handler(c)
			return
		}

		log.Warn().Strs("required", required).Msg("User does not have any of the required permissions")
		helpers.APIErrorResponse(c, http.StatusForbidden, "User does not have any of the required permissions")
		c.Abort()
	}
}

func contains(s []string, e string) bool {
	e = strings.ToLower(e)
	for _, a := range s {
		if strings.ToLower(a) == e {
			return true
		}
	}
	return false
}

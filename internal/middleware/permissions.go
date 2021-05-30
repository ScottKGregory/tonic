package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/api"
	tonicErrors "github.com/scottkgregory/tonic/internal/api/errors"
	"github.com/scottkgregory/tonic/internal/dependencies"
	"github.com/scottkgregory/tonic/internal/services"
)

func HasAny(required ...string) gin.HandlerFunc {
	if valid, messages := services.ValidatePermissions(required...); !valid {
		panic(tonicErrors.NewValidationError(messages))
	}

	return func(c *gin.Context) {
		user, ok := dependencies.GetUser(c)
		if !ok {
			api.ForbiddenResponse(c, tonicErrors.NewForbiddenError(required...))
			c.Abort()
			return
		}

		perms := formatPerms(user.Permissions)
		if contains(c, perms, required...) {
			c.Next()
			return
		}

		api.ForbiddenResponse(c, tonicErrors.NewForbiddenError(required...))
		c.Abort()
	}
}

func HasAll(required ...string) gin.HandlerFunc {
	if valid, messages := services.ValidatePermissions(required...); !valid {
		panic(tonicErrors.NewValidationError(messages))
	}

	return func(c *gin.Context) {
		user, ok := dependencies.GetUser(c)
		if !ok {
			api.ForbiddenResponse(c, tonicErrors.NewForbiddenError(required...))
			c.Abort()
			return
		}

		v := 0
		perms := formatPerms(user.Permissions)
		for _, r := range required {
			if contains(c, perms, r) {
				v += 1
			}
		}
		if v == len(required) {
			c.Next()
			return
		}

		api.ForbiddenResponse(c, tonicErrors.NewForbiddenError(required...))
		c.Abort()
	}
}

func formatPerms(in []string) []string {
	perms := []string{}
	for _, x := range in {
		perms = append(perms, strings.ToLower(x))
	}

	return perms
}

func contains(c *gin.Context, perms []string, required ...string) bool {
	for _, r := range required {
		rs := strings.Split(r, ":")
		for _, y := range perms {
			ys := strings.Split(y, ":")
			v := 0

			for i, p := range rs {
				if i == len(rs)-1 && p != "*" {
					p = c.Param(p)
				}

				if strings.ToLower(p) == ys[i] || ys[i] == "*" {
					v += 1
				}

				if v == 3 {
					return true
				}
			}
		}
	}

	return false
}

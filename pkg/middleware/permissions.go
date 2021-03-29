package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/services"
)

func HasAny(backend backends.Backend, required ...string) gin.HandlerFunc {
	for _, r := range required {
		if len(strings.Split(r, ":")) != 3 {
			panic(fmt.Errorf("permission %s not valid, must have three parts", r))
		}
	}

	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, backend)

		user, err := userService.GetUser(c.GetString(constants.SubjectKey))
		if err != nil {
			api.ForbiddenResponse(c)
			c.Abort()
		}

		perms := formatPerms(user.Permissions)
		if contains(perms, required...) {
			c.Next()
			return
		}

		api.ForbiddenResponse(c)
		c.Abort()
	}
}

func HasAll(backend backends.Backend, required ...string) gin.HandlerFunc {
	for _, r := range required {
		if len(strings.Split(r, ":")) != 3 {
			panic(fmt.Errorf("permission %s not valid, must have three parts", r))
		}
	}

	return func(c *gin.Context) {
		log := dependencies.GetLogger(c)
		userService := services.NewUserService(log, backend)

		user, err := userService.GetUser(c.GetString(constants.SubjectKey))
		if err != nil {
			api.ForbiddenResponse(c)
			c.Abort()
		}

		v := 0
		perms := formatPerms(user.Permissions)
		for _, r := range required {
			if contains(perms, r) {
				v += 1
			}
		}
		if v == len(required) {
			c.Next()
			return
		}

		api.ForbiddenResponse(c)
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

func contains(perms []string, required ...string) bool {
	for _, r := range required {
		r = strings.ToLower(r)
		rs := strings.Split(r, ":")
		for _, y := range perms {
			ys := strings.Split(y, ":")
			v := 0

			for i, p := range rs {
				if p == ys[i] || ys[i] == "*" {
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

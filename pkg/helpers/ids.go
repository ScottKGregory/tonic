package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/dependencies"
	"github.com/scottkgregory/tonic/pkg/models"
)

// GetLogger returns a zerolog logger with context
func GetLogger(c ...*gin.Context) *zerolog.Logger {
	return dependencies.GetLogger(c...)
}

// GetUser returns the current user from context
func GetUser(c *gin.Context) (user *models.User, ok bool) {
	return dependencies.GetUser(c)
}

// IDPath will add the standard ID param to the given path
func IDPath(path ...string) string {
	p := ""
	if len(path) > 0 {
		p = path[0]
	}

	if strings.HasSuffix(p, ":") {
		return p + constants.IDParam
	}

	return strings.TrimSuffix(p, "/") + "/:" + constants.IDParam
}

// GetID gets the id from the path, setup using tonic.IDPath
func GetID(c *gin.Context) (id string) {
	return c.Param(constants.IDParam)
}

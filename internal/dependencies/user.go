package dependencies

import (
	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/constants"
	"github.com/scottkgregory/tonic/pkg/models"
)

func GetUser(c *gin.Context) (user *models.User, ok bool) {
	u, ok := c.Get(constants.UserKey)
	if !ok {
		return nil, false
	} else {
		user, ok = u.(*models.User)
		if !ok {
			return nil, false
		}
	}

	return user, true
}

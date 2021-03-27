package backends

import (
	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/pkg/models"
)

type Backend interface {
	GetUserByOIDCSubject(log *zerolog.Logger, sub string) (user *models.User, err error)
	GetUserByID(log *zerolog.Logger, id string) (user *models.User, err error)
	ListUsers(log *zerolog.Logger) (users *[]models.User, err error)
	SaveUser(log *zerolog.Logger, user *models.User) (err error)
	Ping(log *zerolog.Logger) error
}

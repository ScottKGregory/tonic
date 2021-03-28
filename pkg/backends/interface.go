package backends

import "github.com/scottkgregory/tonic/pkg/models"

type Backend interface {
	CreateUser(in *models.User) (out *models.User, err error)
	UpdateUser(in *models.User) (out *models.User, err error)
	GetUser(subject string) (out *models.User, err error)
	ListUsers() (out *[]models.User, err error)
	Ping() error
}

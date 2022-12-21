package backends

import (
	"context"

	"github.com/scottkgregory/tonic/pkg/models"
)

type Backend interface {
	CreateUser(context.Context, *models.User) (out *models.User, err error)
	UpdateUser(context.Context, *models.User) (out *models.User, err error)
	GetUser(ctx context.Context, subject string) (out *models.User, err error)
	ListUsers(context.Context) (out []*models.User, err error)
	Ping(context.Context) error
}

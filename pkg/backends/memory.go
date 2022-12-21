package backends

import (
	"context"

	"github.com/scottkgregory/tonic/pkg/models"
)

type Memory struct {
	options *models.BackendOptions
}

var _ Backend = Memory{}

var users []*models.User

func NewMemoryBackend(options *models.BackendOptions) *Memory {
	return &Memory{options}
}

func (m Memory) CreateUser(ctx context.Context, in *models.User) (out *models.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == in.Claims.Subject {
			*u = *in
			return u, nil
		}
	}

	users = append(users, in)

	return in, err
}

func (m Memory) UpdateUser(ctx context.Context, in *models.User) (out *models.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == in.Claims.Subject {
			*u = *in
			return u, nil
		}
	}

	return in, err
}

func (m Memory) GetUser(ctx context.Context, subject string) (out *models.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == subject {
			return u, nil
		}
	}

	return nil, nil
}

func (m Memory) ListUsers(ctx context.Context) (out []*models.User, err error) {
	return users, nil
}

func (m Memory) Ping(ctx context.Context) error {
	return nil
}

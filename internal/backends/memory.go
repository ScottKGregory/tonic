package backends

import (
	"github.com/scottkgregory/tonic/internal/models"
	pkgModels "github.com/scottkgregory/tonic/pkg/models"
)

type Memory struct {
	options *models.BackendOptions
}

var _ Backend = Memory{}

var users []*pkgModels.User

func NewMemoryBackend(options *models.BackendOptions) *Memory {
	return &Memory{options}
}

func (m Memory) CreateUser(in *pkgModels.User) (out *pkgModels.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == in.Claims.Subject {
			*u = *in
			return u, nil
		}
	}

	users = append(users, in)

	return in, err
}

func (m Memory) UpdateUser(in *pkgModels.User) (out *pkgModels.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == in.Claims.Subject {
			*u = *in
			return u, nil
		}
	}

	return in, err
}

func (m Memory) GetUser(subject string) (out *pkgModels.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == subject {
			return u, nil
		}
	}

	return nil, nil
}

func (m Memory) ListUsers() (out []*pkgModels.User, err error) {
	return users, nil
}

func (m Memory) Ping() error {
	return nil
}

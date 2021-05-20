package backends

import "github.com/scottkgregory/tonic/pkg/models"

type Memory struct {
	options *models.Backend
}

var _ Backend = Memory{}

var users []*models.User

func NewMemoryBackend(options *models.Backend) *Memory {
	return &Memory{options}
}

func (m Memory) CreateUser(in *models.User) (out *models.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == in.Claims.Subject {
			*u = *in
			return u, nil
		}
	}

	users = append(users, in)

	return in, err
}

func (m Memory) UpdateUser(in *models.User) (out *models.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == in.Claims.Subject {
			*u = *in
			return u, nil
		}
	}

	return in, err
}

func (m Memory) GetUser(subject string) (out *models.User, err error) {
	for _, u := range users {
		if u.Claims.Subject == subject {
			return u, nil
		}
	}

	return nil, nil
}

func (m Memory) ListUsers() (out []*models.User, err error) {
	return users, nil
}

func (m Memory) Ping() error {
	return nil
}

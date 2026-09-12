package adapters

import (
	"slices"
	"strconv"

	"codeberg.org/MrTomate/web/internal/core"
)

type InMemoryUserRepository struct {
	users []*core.User
}

func (r *InMemoryUserRepository) CreateUser(u *core.User) (*core.User, error) {
	id := len(r.users)
	u.ID = strconv.Itoa(id + 1)

	i :=  slices.IndexFunc(r.users,func(e *core.User) bool {
		return e.Email == u.Email
	})
	
	if i != -1 {
		return nil, core.ErrEmailAlreadyUsed
	}
	
	r.users = append(r.users, u)
	return u, nil
}
 
func (r *InMemoryUserRepository) UpdateUser(u *core.User) (*core.User, error) {
	for i, user := range r.users {
		if user.Email == u.Email && user.ID == u.ID {
			r.users[i] = u
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}

func (r *InMemoryUserRepository) FindById(id string) (*core.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*core.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make([]*core.User, 0),
	}
}

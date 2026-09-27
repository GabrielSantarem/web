package adapters

import (
	"errors"
	"slices"
	"strconv"
	"sync"

	"codeberg.org/MrTomate/web/internal/core"
)

// InMemoryUserRepository implementa a porta core.UserRepository
// reutilizando a base generica MemoryRepository[*core.User].
type InMemoryUserRepository struct {
	*MemoryRepository[*core.User]
	mu    sync.RWMutex
	users []*core.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		MemoryRepository: NewMemoryRepository[*core.User](),
		users:            make([]*core.User, 0),
	}
}

func (r *InMemoryUserRepository) Create(u *core.User) (*core.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i := slices.IndexFunc(r.users, func(e *core.User) bool {
		return e.Email == u.Email
	})

	if i != -1 {
		return nil, core.ErrEmailAlreadyUsed
	}

	if u.ID == "" {
		id := len(r.users)
		u.ID = strconv.Itoa(id + 1)
	}

	if _, err := r.MemoryRepository.Create(u); err != nil {
		return nil, err
	}

	r.users = append(r.users, u)
	return u, nil
}

func (r *InMemoryUserRepository) CreateUser(u *core.User) (*core.User, error) {
	return r.Create(u)
}

func (r *InMemoryUserRepository) Update(u *core.User) (*core.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	updated, err := r.MemoryRepository.Update(u)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}

	for i, user := range r.users {
		if user.ID == u.ID {
			r.users[i] = u
			break
		}
	}

	return updated, nil
}

func (r *InMemoryUserRepository) UpdateUser(u *core.User) (*core.User, error) {
	return r.Update(u)
}

func (r *InMemoryUserRepository) FindByID(id string) (*core.User, error) {
	u, err := r.MemoryRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *InMemoryUserRepository) FindById(id string) (*core.User, error) {
	return r.FindByID(id)
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*core.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}

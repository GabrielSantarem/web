package core_test

import (
	"errors"
	"testing"

	"codeberg.org/MrTomate/web/internal/core"
)

type mockUserRepository struct {
	users       map[string]*core.User
	createErr   error
	updateErr   error
	findByIDErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*core.User),
	}
}

func (m *mockUserRepository) CreateUser(u *core.User) (*core.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	for _, existing := range m.users {
		if existing.Email == u.Email {
			return nil, core.ErrEmailAlreadyUsed
		}
	}
	u.ID = "id-123"
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserRepository) UpdateUser(u *core.User) (*core.User, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	if _, ok := m.users[u.ID]; !ok {
		return nil, core.ErrUserNotFound
	}
	m.users[u.ID] = u
	return u, nil
}

func (m *mockUserRepository) FindById(id string) (*core.User, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	u, ok := m.users[id]
	if !ok {
		return nil, core.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepository) FindByEmail(email string) (*core.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}

func TestUserService(t *testing.T) {
	t.Run("NewUser inicializa corretamente", func(t *testing.T) {
		u := core.NewUser("test@example.com", "Test User")
		if u.Email != "test@example.com" || u.Name != "Test User" {
			t.Errorf("dados incorretos: %+v", u)
		}
		if u.IsActive {
			t.Error("novo usuario deve iniciar inativo")
		}
	})

	t.Run("CreateUser com dados validos cria usuario", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		u := core.NewUser("novo@example.com", "Novo Usuario")
		created, err := service.CreateUser(u)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if created.ID != "id-123" {
			t.Errorf("esperava id-123, veio %s", created.ID)
		}
		if created.GetID() != "id-123" {
			t.Errorf("GetID() diferente do ID: %s", created.GetID())
		}
	})

	t.Run("CreateUser rejeita email vazio", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		u := core.NewUser("", "Nome Valido")
		_, err := service.CreateUser(u)
		if !errors.Is(err, core.ErrInvalidInput) {
			t.Fatalf("esperava ErrInvalidInput, veio: %v", err)
		}
	})

	t.Run("CreateUser rejeita nome vazio", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		u := core.NewUser("email@valido.com", "")
		_, err := service.CreateUser(u)
		if !errors.Is(err, core.ErrInvalidInput) {
			t.Fatalf("esperava ErrInvalidInput, veio: %v", err)
		}
	})

	t.Run("GetUser busca usuario com sucesso", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		created, _ := service.CreateUser(core.NewUser("busca@example.com", "Busca"))
		found, err := service.GetUser(created.ID)
		if err != nil {
			t.Fatalf("erro ao buscar usuario: %v", err)
		}
		if found.Name != "Busca" {
			t.Errorf("esperava Busca, veio %s", found.Name)
		}
	})

	t.Run("GetUser retorna erro quando usuario nao existe", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		_, err := service.GetUser("nao-existe")
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("esperava ErrUserNotFound, veio: %v", err)
		}
	})

	t.Run("ActivateUser ativa usuario inativo com sucesso", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		created, _ := service.CreateUser(core.NewUser("ativo@example.com", "Ativar"))
		if created.IsActive {
			t.Fatal("esperava que o usuario estivesse inativo")
		}

		err := service.ActivateUser(created.ID)
		if err != nil {
			t.Fatalf("erro ao ativar usuario: %v", err)
		}

		updated, _ := service.GetUser(created.ID)
		if !updated.IsActive {
			t.Error("usuario deveria estar ativo")
		}
	})

	t.Run("ActivateUser falha se usuario ja estiver ativo", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		created, _ := service.CreateUser(core.NewUser("ja_ativo@example.com", "Ja Ativo"))
		_ = service.ActivateUser(created.ID)

		err := service.ActivateUser(created.ID)
		if !errors.Is(err, core.ErrUserAlreadyActive) {
			t.Fatalf("esperava ErrUserAlreadyActive, veio: %v", err)
		}
	})

	t.Run("ActivateUser falha se usuario nao for encontrado", func(t *testing.T) {
		mockRepo := newMockUserRepository()
		service := core.NewUserService(mockRepo)

		err := service.ActivateUser("inexistente")
		if !errors.Is(err, core.ErrUserNotFound) {
			t.Fatalf("esperava ErrUserNotFound, veio: %v", err)
		}
	})
}

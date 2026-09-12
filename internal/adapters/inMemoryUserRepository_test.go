package adapters

import (
	"testing"

	"codeberg.org/MrTomate/web/internal/core"
)

func TestInMemoryUserRepository(t *testing.T) {
	repo := NewInMemoryUserRepository()
	if repo == nil {
		t.Fatal("NewInMemoryUserRepository returned nil")
	}
	if len(repo.users) != 0 {
		t.Fatal("users slice should be empty")
	}
}

func TestInMemoryUserRepository_AddUser(t *testing.T) {
	t.Run("salva um user no repo", func(t *testing.T) {
		repo := NewInMemoryUserRepository()
		user := &core.User{
			ID:   "1",
			Name: "test",
		}
		_, err := repo.CreateUser(user)
		if err != nil {
			t.Fatal(err)
		}
		if len(repo.users) != 1 {
			t.Fatal("users slice should have 1 element")
		}
	})

	t.Run("Email duplicado retorna um error", func(t *testing.T) {
	
	})
}

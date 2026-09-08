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
	repo := NewInMemoryUserRepository()
	user := &core.User{
		ID:   "1",
		Name: "test",
	}
	_, err := repo.Save(user)
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.users) != 1 {
		t.Fatal("users slice should have 1 element")
	}
}

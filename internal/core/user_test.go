package core_test

import (
	"testing"

	"codeberg.org/MrTomate/web/internal/adapters"
	"codeberg.org/MrTomate/web/internal/core"
)

func TestNewEmptyUser(t *testing.T) {
	u := core.NewEmptyUser()
	if u == nil {
		t.Errorf("NewEmptyUser() returned nil")
	}
	if u.ID != "" {
		t.Errorf("ID should be empty")
	}
	if u.Email != "" {
		t.Errorf("Email should be empty")
	}
	if u.Name != "" {
		t.Errorf("Name should be empty")
	}
}

func TestCreateUser_WithEmptyUser_ReturnsError(t *testing.T) {
	repo := adapters.NewInMemoryUserRepository()
	srv := core.NewUserService(repo)
	u := core.NewEmptyUser()

	user, err := srv.CreateUser(u)
	if err == nil {
		t.Errorf("CreateUser() should return an error for an empty user")
	}
	if user != nil {
		t.Errorf("CreateUser() should return nil user for invalid input")
	}
}

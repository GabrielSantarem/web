package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codeberg.org/MrTomate/web/handler"
	"codeberg.org/MrTomate/web/internal/adapters"
	"codeberg.org/MrTomate/web/internal/core"
)

func TestHandleCreateUser(t *testing.T) {
	// factory: cada subtest ganha um handler novo com repo limpo,
	// sem estado compartilhado entre testes
	newHandler := func() http.HandlerFunc {
		repo := adapters.NewInMemoryUserRepository()
		return handler.HandleCreateUser(core.NewUserService(repo))
	}

	t.Run("user valido retorna 201 com o user criado", func(t *testing.T) {
		body := strings.NewReader(`{"name":"tomate","email":"tomate@test.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/user/create", body)
		rec := httptest.NewRecorder()

		newHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status: want %d, got %d", http.StatusCreated, rec.Code)
		}

		var got core.User
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decoding response body: %v", err)
		}

		if got.Name != "tomate" {
			t.Errorf("name: want %q, got %q", "tomate", got.Name)
		}
		if got.Email != "tomate@test.com" {
			t.Errorf("email: want %q, got %q", "tomate@test.com", got.Email)
		}
		if got.ID == "" {
			t.Error("expected non-empty ID")
		}
	})

	t.Run("email invalido retorna 400 com erro no body", func(t *testing.T) {
		body := strings.NewReader(`{"name":"tomate","email":"tomate#test.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/user/create", body)
		rec := httptest.NewRecorder()

		newHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: want %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var got struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decoding response body: %v", err)
		}
		if got.Error == "" {
			t.Error("expected non-empty error message")
		}
	})

	t.Run("json malformado retorna 400", func(t *testing.T) {
		body := strings.NewReader(`{invalido`)
		req := httptest.NewRequest(http.MethodPost, "/user/create", body)
		rec := httptest.NewRecorder()

		newHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status: want %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

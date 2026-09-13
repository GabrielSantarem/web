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

		var got handler.ErrorResponse
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decoding response body: %v", err)
		}
		if got.Message == "" && got.Error() == "" {
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

	t.Run("json invalidos devem retorna um erro", func(t *testing.T) {
		body := strings.NewReader(`{"name":"To","email":"tomate#teste"}`)
		req := httptest.NewRequest(http.MethodPost, "/user/create", body)
		rec := httptest.NewRecorder()

		var errBody handler.ErrorResponse

		h := newHandler()
		h.ServeHTTP(rec, req)

		if err := json.NewDecoder(rec.Body).Decode(&errBody); err != nil {
			t.Fatal("parse body response failure")
		}

		if errBody.Message != handler.ErrValidationFailure.Error() {
			t.Fatalf("validation: expected %s, got %s", handler.ErrValidationFailure.Error(), errBody.Message)
		}
	})

	t.Run("email duplicado retorna 409 conflict", func(t *testing.T) {
		payload := `{"name":"tomate","email":"tomate@gmail.com"}`
		h := newHandler()

		// Primeiro insert
		req1 := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(payload))
		rec1 := httptest.NewRecorder()
		h.ServeHTTP(rec1, req1)

		if rec1.Code != http.StatusCreated {
			t.Fatalf("first insert status: want %d, got %d", http.StatusCreated, rec1.Code)
		}

		// Segundo insert (instanciando uma nova requisição com novo buffer)
		req2 := httptest.NewRequest(http.MethodPost, "/user/create", strings.NewReader(payload))
		rec2 := httptest.NewRecorder()
		h.ServeHTTP(rec2, req2)

		if rec2.Code != http.StatusConflict {
			t.Fatalf("second insert status: want %d, got %d", http.StatusConflict, rec2.Code)
		}
	})
}

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/MrTomate/web/types"
)

func TestHandleHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h := HandleHello()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("espera status:%d ,veio %d", http.StatusAccepted, rec.Code)
	}

	gotContentType := rec.Header().Get("Content-Type")

	if gotContentType != "application/json" {
		t.Fatalf("esperava Content-Type application/json veio %q", gotContentType)
	}

	var got types.Message

	err := json.NewDecoder(rec.Body).Decode(&got)
	if err != nil {
		t.Fatalf("erro ao decodificar json: %v", err)
	}
}

func TestHandleMe(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/me/Tomate", nil)
	req.SetPathValue("Name", "Tomate")

	rec := httptest.NewRecorder()

	h := HandleMe()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava status %d, veio %d", http.StatusOK, rec.Code)
	}

	var got types.Me

	err := json.NewDecoder(rec.Body).Decode(&got)

	if err != nil {
		t.Fatalf("erro ao decodificar json: %v", err)
	}

	if got.Name != "Tomate" {
		t.Fatalf("esperava Name = Tomate, veio %q", got.Name)
	}
}

func TestHandleNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/random-bulshit", nil)
	rec := httptest.NewRecorder()

	h := HandleNotFound()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("esperava status %d, veio %d", http.StatusSeeOther, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/" {
		t.Fatalf("espereva redirect pra '/' veio %q", location)
	}
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"messaging-application/servers/gateway/models/users"
	"messaging-application/servers/gateway/sessions"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var ctx *Context

const secret = "c2VjcmV0"

func newContext() *http.ServeMux {
	ctx = &Context{
		secret,
		sessions.NewMemoryStore(),
		users.NewStubStore(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/{UserID}", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/{SessionID}", ctx.SpecificSessionHandler)
	return mux
}

func TestNewUser(t *testing.T) {
	handler := newContext()

	jsonData, err := json.Marshal(&users.NewUser{
		Password:  "password343",
		Email:     "valid_email@example.com",
		FirstName: "Jon",
		LastName:  "Doe",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	fmt.Print(rr.Body.String())

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}
	if !strings.HasPrefix(rr.Header().Get("Authorization"), "Bearer ") {
		t.Error("Expected Authorization header to be set")
	}
	if !strings.Contains(rr.Body.String(), `"id":1`) {
		t.Error("Expected user ID to be returned")
	}
}

func TestNewUserErrors(t *testing.T) {
	handler := newContext()

	// Wrong Method
	req, err := http.NewRequest(http.MethodGet, "/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, status)
	}

	// Wrong Content Type
	req, err = http.NewRequest(http.MethodPost, "/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status %d, got %d", http.StatusUnsupportedMediaType, status)
	}

	// Invalid Body
	bodies := []*users.NewUser{
		{
			Password: "",
		},
		{
			Password: "not_empty",
			Email:    "invalid_email",
		},
		{
			Password:  "not_empty",
			Email:     "valid_email@example.com",
			FirstName: "345346invalid_name",
		},
		{
			Password:  "not_empty",
			Email:     "valid_email@example.com",
			FirstName: "Validname",
			LastName:  "8y78fa6sf8invalid_name",
		},
	}

	for _, body := range bodies {
		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
		}
	}
}

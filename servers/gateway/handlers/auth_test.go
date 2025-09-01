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

func TestSpecificUser(t *testing.T) {
	handler := newContext()

	// create new user
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
	if rr.Code != http.StatusCreated {
		t.Fatal(rr.Body.String())
	}

	authHeader := rr.Header().Get("Authorization")

	// retrieve the user
	req, err = http.NewRequest(http.MethodGet, "/v1/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	if !strings.Contains(rr.Body.String(), `"id":1`) {
		t.Error("Expected user ID to be returned")
	}
	if !strings.Contains(rr.Body.String(), `"firstName":"Jon"`) {
		t.Error("Expected first name to be returned")
	}
	if !strings.Contains(rr.Body.String(), `"lastName":"Doe"`) {
		t.Error("Expected last name to be returned")
	}

	// modify user
	jsonData, err = json.Marshal(&users.Updates{
		Email:     "new_email@example.com",
		FirstName: "NewFirstName",
		LastName:  "NewLastName",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest(http.MethodPatch, "/v1/users/1", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
	if !strings.Contains(rr.Body.String(), `"firstName":"NewFirstName"`) {
		t.Error("Expected first name to be returned")
	}
	if !strings.Contains(rr.Body.String(), `"lastName":"NewLastName"`) {
		t.Error("Expected last name to be returned")
	}
}

func TestSpecificUserErrors(t *testing.T) {
	handler := newContext()

	// create new user
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
	if rr.Code != http.StatusCreated {
		t.Fatal(rr.Body.String())
	}

	authHeader := rr.Header().Get("Authorization")

	// Missing auth header
	req, err = http.NewRequest(http.MethodGet, "/v1/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
	}

	// Wrong method
	req, err = http.NewRequest(http.MethodPost, "/v1/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, status)
	}

	// Non-existent user
	req, err = http.NewRequest(http.MethodGet, "/v1/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}

	// Forbidden user modification
	req, err = http.NewRequest(http.MethodPatch, "/v1/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, status)
	}

	// Wrong content type
	req, err = http.NewRequest(http.MethodPatch, "/v1/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "text/plain")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status %d, got %d", http.StatusUnsupportedMediaType, status)
	}
}

func TestSessionsFlow(t *testing.T) {
	handler := newContext()

	// Create a new user
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
	if rr.Code != http.StatusCreated {
		t.Fatal(rr.Body.String())
	}

	authHeader := rr.Header().Get("Authorization")

	// Delete session
	req, err = http.NewRequest(http.MethodDelete, "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, status)
	}

	// Create new session
	jsonData, err = json.Marshal(&users.Credentials{
		Password: "password343",
		Email:    "valid_email@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest(http.MethodPost, "/v1/sessions", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	authHeader = rr.Header().Get("Authorization")

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	// Retrieve user
	req, err = http.NewRequest(http.MethodGet, "/v1/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestSessionsErrors(t *testing.T) {
	handler := newContext()

	// Create a new user
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
	if rr.Code != http.StatusCreated {
		t.Fatal(rr.Body.String())
	}

	authHeader := rr.Header().Get("Authorization")

	// Wrong content type for create session
	req, err = http.NewRequest(http.MethodPost, "/v1/sessions", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status %d, got %d", http.StatusUnsupportedMediaType, status)
	}

	// Wrong method for create session
	req, err = http.NewRequest(http.MethodPut, "/v1/sessions", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, status)
	}

	// Wrong method for delete session
	req, err = http.NewRequest(http.MethodPut, "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, status)
	}

	// Forbidden session deletion
	req, err = http.NewRequest(http.MethodDelete, "/v1/sessions/999", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", authHeader)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, status)
	}

	// Login to non-existent account
	jsonData, err = json.Marshal(&users.Credentials{
		Password: "password343",
		Email:    "non_existent_email@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest(http.MethodPost, "/v1/sessions", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
	}

	// Login with incorrect password
	jsonData, err = json.Marshal(&users.Credentials{
		Password: "wrong_password",
		Email:    "valid_email@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest(http.MethodPost, "/v1/sessions", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, status)
	}
}

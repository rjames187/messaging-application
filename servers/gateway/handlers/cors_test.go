package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created"))
}

func TestCORSHandlerNoPreFlight(t *testing.T) {
	handler := NewCORSHandler(http.HandlerFunc(mockHandler))
	req := httptest.NewRequest(http.MethodPost, "https://example.com", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Result().StatusCode)
	}
}

func TestCORSHandlerWithPreFlight(t *testing.T) {
	handler := NewCORSHandler(http.HandlerFunc(mockHandler))
	req := httptest.NewRequest(http.MethodOptions, "https://example.com", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Result().StatusCode)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin header to be '*', got '%s'", rr.Header().Get("Access-Control-Allow-Origin"))
	}
	if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, PATCH, DELETE, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods header to be 'GET, POST, PUT, PATCH, DELETE, OPTIONS', got '%s'", rr.Header().Get("Access-Control-Allow-Methods"))
	}
	if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("Expected Access-Control-Allow-Headers header to be 'Content-Type, Authorization', got '%s'", rr.Header().Get("Access-Control-Allow-Headers"))
	}
	if rr.Header().Get("Access-Control-Expose-Headers") != "Authorization" {
		t.Errorf("Expected Access-Control-Expose-Headers header to be 'Authorization', got '%s'", rr.Header().Get("Access-Control-Expose-Headers"))
	}
	if rr.Header().Get("Access-Control-Max-Age") != oneDay {
		t.Errorf("Expected Access-Control-Max-Age header to be '%s', got '%s'", oneDay, rr.Header().Get("Access-Control-Max-Age"))
	}
}

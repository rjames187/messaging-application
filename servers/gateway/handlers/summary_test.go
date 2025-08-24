package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchHTML(t *testing.T) {
	// Test successful HTML fetch
	t.Run("successful fetch", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<html><head><title>Test Title</title></head></html>`))
		}))
		defer server.Close()

		metadata, err := fetchHTML(server.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if metadata == nil {
			t.Fatal("expected metadata, got nil")
		}
		if metadata.Title != "Test Title" {
			t.Errorf("expected title 'Test Title', got '%s'", metadata.Title)
		}
	})

	// Test invalid URL
	t.Run("invalid URL", func(t *testing.T) {
		_, err := fetchHTML("invalid-url")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
		if !strings.HasPrefix(err.Error(), "error fetching html:") {
			t.Errorf("expected error to start with 'error fetching html:', got '%s'", err.Error())
		}
	})

	// Test unreachable server
	t.Run("unreachable server", func(t *testing.T) {
		_, err := fetchHTML("http://localhost:99999")
		if err == nil {
			t.Fatal("expected error for unreachable server")
		}
		if !strings.HasPrefix(err.Error(), "error fetching html:") {
			t.Errorf("expected error to start with 'error fetching html:', got '%s'", err.Error())
		}
	})
}

func TestExtractSummary(t *testing.T) {
	t.Run("extract title", func(t *testing.T) {
		htmlContent := `<html><head><title>Test Page Title</title></head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if metadata.Title != "Test Page Title" {
			t.Errorf("expected title 'Test Page Title', got '%s'", metadata.Title)
		}
	})

	t.Run("extract og meta tags", func(t *testing.T) {
		htmlContent := `<html><head>
			<meta property="og:title" content="OG Title">
			<meta property="og:description" content="OG Description">
			<meta property="og:url" content="https://example.com">
			<meta property="og:type" content="website">
			<meta property="og:site_name" content="Example Site">
		</head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if metadata.Title != "OG Title" {
			t.Errorf("expected title 'OG Title', got '%s'", metadata.Title)
		}
		if metadata.Description != "OG Description" {
			t.Errorf("expected description 'OG Description', got '%s'", metadata.Description)
		}
		if metadata.URL != "https://example.com" {
			t.Errorf("expected URL 'https://example.com', got '%s'", metadata.URL)
		}
		if metadata.Type != "website" {
			t.Errorf("expected type 'website', got '%s'", metadata.Type)
		}
		if metadata.SiteName != "Example Site" {
			t.Errorf("expected site name 'Example Site', got '%s'", metadata.SiteName)
		}
	})

	t.Run("extract twitter meta tags", func(t *testing.T) {
		htmlContent := `<html><head>
			<meta property="twitter:card" content="summary">
			<meta property="twitter:title" content="Twitter Title">
			<meta property="twitter:description" content="Twitter Description">
		</head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if metadata.Type != "summary" {
			t.Errorf("expected type 'summary', got '%s'", metadata.Type)
		}
		if metadata.Title != "Twitter Title" {
			t.Errorf("expected title 'Twitter Title', got '%s'", metadata.Title)
		}
		if metadata.Description != "Twitter Description" {
			t.Errorf("expected description 'Twitter Description', got '%s'", metadata.Description)
		}
	})

	t.Run("extract standard meta tags", func(t *testing.T) {
		htmlContent := `<html><head>
			<meta name="description" content="Standard Description">
			<meta name="author" content="John Doe">
			<meta name="keywords" content="go, testing, html">
		</head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if metadata.Description != "Standard Description" {
			t.Errorf("expected description 'Standard Description', got '%s'", metadata.Description)
		}
		if metadata.Author != "John Doe" {
			t.Errorf("expected author 'John Doe', got '%s'", metadata.Author)
		}
		if len(metadata.Keywords) != 3 {
			t.Errorf("expected 3 keywords, got %d", len(metadata.Keywords))
		}
		expectedKeywords := []string{"go", "testing", "html"}
		for i, keyword := range expectedKeywords {
			if metadata.Keywords[i] != keyword {
				t.Errorf("expected keyword '%s', got '%s'", keyword, metadata.Keywords[i])
			}
		}
	})

	t.Run("extract og images", func(t *testing.T) {
		htmlContent := `<html><head>
			<meta property="og:image" content="https://example.com/image.jpg">
			<meta property="og:image:width" content="800">
			<meta property="og:image:height" content="600">
			<meta property="og:image:type" content="image/jpeg">
			<meta property="og:image:alt" content="Example Image">
		</head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if len(metadata.Images) != 1 {
			t.Fatalf("expected 1 image, got %d", len(metadata.Images))
		}
		image := metadata.Images[0]
		if image.URL != "https://example.com/image.jpg" {
			t.Errorf("expected image URL 'https://example.com/image.jpg', got '%s'", image.URL)
		}
		if image.Width != 800 {
			t.Errorf("expected image width 800, got %d", image.Width)
		}
		if image.Height != 600 {
			t.Errorf("expected image height 600, got %d", image.Height)
		}
		if image.Type != "image/jpeg" {
			t.Errorf("expected image type 'image/jpeg', got '%s'", image.Type)
		}
		if image.Alt != "Example Image" {
			t.Errorf("expected image alt 'Example Image', got '%s'", image.Alt)
		}
	})

	t.Run("extract link icon", func(t *testing.T) {
		htmlContent := `<html><head>
			<link rel="icon" href="/favicon.ico" type="image/x-icon" sizes="16x16">
		</head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if metadata.Icon.URL != "/favicon.ico" {
			t.Errorf("expected icon URL '/favicon.ico', got '%s'", metadata.Icon.URL)
		}
		if metadata.Icon.Type != "image/x-icon" {
			t.Errorf("expected icon type 'image/x-icon', got '%s'", metadata.Icon.Type)
		}
		if metadata.Icon.Width != 16 {
			t.Errorf("expected icon width 16, got %d", metadata.Icon.Width)
		}
		if metadata.Icon.Height != 16 {
			t.Errorf("expected icon height 16, got %d", metadata.Icon.Height)
		}
	})

	t.Run("priority handling", func(t *testing.T) {
		htmlContent := `<html><head>
			<title>HTML Title</title>
			<meta name="description" content="Standard Description">
			<meta property="twitter:title" content="Twitter Title">
			<meta property="twitter:description" content="Twitter Description">
			<meta property="og:title" content="OG Title">
			<meta property="og:description" content="OG Description">
		</head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		// og:title should override twitter:title and HTML title
		if metadata.Title != "OG Title" {
			t.Errorf("expected title 'OG Title', got '%s'", metadata.Title)
		}
		// og:description should override twitter:description and standard description
		if metadata.Description != "OG Description" {
			t.Errorf("expected description 'OG Description', got '%s'", metadata.Description)
		}
	})

	t.Run("empty HTML", func(t *testing.T) {
		htmlContent := `<html><head></head></html>`
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader(htmlContent)),
		}

		metadata := extractSummary(resp)
		if metadata == nil {
			t.Fatal("expected metadata, got nil")
		}
		if metadata.Title != "" {
			t.Errorf("expected empty title, got '%s'", metadata.Title)
		}
	})
}

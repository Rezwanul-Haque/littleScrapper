package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHTMLContent(t *testing.T) {
	// Create a test server with a sample HTML response
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("<html><head><title>Test Page</title></head><body><p>Sample Text</p></body></html>"))
	}))
	defer server.Close()

	// Test fetching HTML content
	doc, err := getHTMLContent(server.URL)
	if err != nil {
		t.Errorf("Failed to fetch HTML content: %v", err)
	}

	// Verify the title and body text
	title, body := extractTitleAndBody(doc)
	expectedTitle := "Test Page"
	expectedBody := "Sample Text"
	if title != expectedTitle || body != expectedBody {
		t.Errorf("Title and body do not match. Expected: %s, %s, Got: %s, %s", expectedTitle, expectedBody, title, body)
	}
}

// Additional test cases for other functions can be added similarly

func TestSanitizeFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"image.jpg?width=120&dpr=1", "image.jpg"},
		{"image with spaces.png", "image_with_spaces.png"},
		{"image&special$chars.png", "image_special_chars.png"},
	}

	for _, testCase := range testCases {
		result := sanitizeFilename(testCase.input)
		if result != testCase.expected {
			t.Errorf("SanitizeFilename failed for input %s. Expected: %s, Got: %s", testCase.input, testCase.expected, result)
		}
	}
}

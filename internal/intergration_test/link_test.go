package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURLIntegration(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("INSTANCE_ID", "test-instance-123")
	os.Setenv("PORT", "8083")

	// Create config
	config, err := api.NewConfig()
	require.NoError(t, err)

	// Create engine with nil Redis client for testing
	engine := api.NewEngine(config, nil)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		engine.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Test request
	reqBody := map[string]interface{}{
		"url": "https://example.com",
		"exp": 604800,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Make request
	resp, err := http.Post(server.URL+"/v1/links/shorten", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEmpty(t, response["code"])
	assert.Equal(t, "Shorten URL generated successfully!", response["message"])
	assert.Len(t, response["code"], 7)
}

func TestShortenURLIntegrationInvalidRequest(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("INSTANCE_ID", "test-instance-123")
	os.Setenv("PORT", "8084")

	// Create config
	config, err := api.NewConfig()
	require.NoError(t, err)

	// Create engine with nil Redis client for testing
	engine := api.NewEngine(config, nil)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		engine.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Test request with missing fields
	reqBody := `{"url": "https://example.com"}` // Missing exp field

	// Make request
	resp, err := http.Post(server.URL+"/v1/links/shorten", "application/json", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShortenURLIntegrationInvalidJSON(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("INSTANCE_ID", "test-instance-123")
	os.Setenv("PORT", "8085")

	// Create config
	config, err := api.NewConfig()
	require.NoError(t, err)

	// Create engine with nil Redis client for testing
	engine := api.NewEngine(config, nil)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		engine.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Test request with invalid JSON
	reqBody := `{"url": "invalid json`

	// Make request
	resp, err := http.Post(server.URL+"/v1/links/shorten", "application/json", bytes.NewBufferString(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

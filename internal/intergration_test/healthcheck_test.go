package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheckIntegration(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("INSTANCE_ID", "test-instance-123")
	os.Setenv("PORT", "8081")

	// Create config
	config, err := api.NewConfig()
	require.NoError(t, err)
	assert.Equal(t, "test-service", config.ServiceName)
	assert.Equal(t, "test-instance-123", config.InstanceID)

	// Create engine with nil Redis client for testing
	engine := api.NewEngine(config, nil)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		engine.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL + "/health-check")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "OK", response["message"])
	assert.Equal(t, "test-service", response["service_name"])
	assert.Equal(t, "test-instance-123", response["instance_id"])
}

func TestHealthCheckIntegrationWithGeneratedUUID(t *testing.T) {
	// Set environment variables without INSTANCE_ID to test UUID generation
	os.Setenv("SERVICE_NAME", "test-service-uuid")
	os.Unsetenv("INSTANCE_ID")
	os.Setenv("PORT", "8082")

	// Create config
	config, err := api.NewConfig()
	require.NoError(t, err)
	assert.Equal(t, "test-service-uuid", config.ServiceName)
	assert.NotEmpty(t, config.InstanceID)
	// Verify it's a valid UUID format (basic check)
	assert.Len(t, config.InstanceID, 36)

	// Create engine with nil Redis client for testing
	engine := api.NewEngine(config, nil)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		engine.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL + "/health-check")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "OK", response["message"])
	assert.Equal(t, "test-service-uuid", response["service_name"])
	assert.Equal(t, config.InstanceID, response["instance_id"])
}

func TestHealthCheckIntegrationDefaultValues(t *testing.T) {
	// Unset all environment variables to test defaults
	os.Unsetenv("SERVICE_NAME")
	os.Unsetenv("INSTANCE_ID")
	os.Unsetenv("PORT")

	// Create config
	config, err := api.NewConfig()
	require.NoError(t, err)
	assert.Equal(t, "health-check-service", config.ServiceName)
	assert.NotEmpty(t, config.InstanceID)
	assert.Equal(t, "8080", config.Port)

	// Create engine with nil Redis client for testing
	engine := api.NewEngine(config, nil)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		engine.ServeHTTP(w, r)
	}))
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL + "/health-check")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "OK", response["message"])
	assert.Equal(t, "health-check-service", response["service_name"])
	assert.Equal(t, config.InstanceID, response["instance_id"])
}

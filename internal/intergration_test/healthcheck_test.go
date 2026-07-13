package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	// not t.Parallel() because of os.Setenv

	testCases := []struct {
		name               string
		setupEnv           func() // setup env for each case
		expectedStatusCode int
		expectedMessage    string
		verifyResponse     func(t *testing.T, body string) // verify response
	}{
		{
			name: "With Set Env",
			setupEnv: func() {
				os.Setenv("SERVICE_NAME", "test-service")
				os.Setenv("INSTANCE_ID", "test-instance-123")
				os.Setenv("APP_PORT", "8081")
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			verifyResponse: func(t *testing.T, body string) {
				var resp map[string]string
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)
				assert.Equal(t, "test-service", resp["service_name"])
				assert.Equal(t, "test-instance-123", resp["instance_id"])
			},
		},
		{
			name: "With Generated UUID",
			setupEnv: func() {
				os.Setenv("SERVICE_NAME", "test-service-uuid")
				os.Unsetenv("INSTANCE_ID") // unset env to generate UUID
				os.Setenv("APP_PORT", "8082")
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			verifyResponse: func(t *testing.T, body string) {
				var resp map[string]string
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)
				assert.Equal(t, "test-service-uuid", resp["service_name"])
				// ensure instance_id has length of UUID (36 characters)
				assert.Len(t, resp["instance_id"], 36) 
			},
		},
		{
			name: "With default env",
			setupEnv: func() {
				// clear env to check default values
				os.Unsetenv("SERVICE_NAME")
				os.Unsetenv("INSTANCE_ID")
				os.Unsetenv("APP_PORT")
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "OK",
			verifyResponse: func(t *testing.T, body string) {
				var resp map[string]string
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)
				// default service name is "health-check-service"
				assert.Equal(t, "health-check-service", resp["service_name"]) 
				assert.Len(t, resp["instance_id"], 36)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup environment
			tc.setupEnv()

			// create config (NewConfig will read from env)
			cfg, err := api.NewConfig()
			require.NoError(t, err)

			// create virtual database miniredis
			mr, err := miniredis.Run()
			require.NoError(t, err)
			defer mr.Close()

			// create client connected to miniredis
			redisClient := redis.NewClient(&redis.Options{
				Addr: mr.Addr(),
			})
			defer redisClient.Close()

			// create engine application
			apiEngine := api.NewEngine(cfg, redisClient)

			// simulate HTTP request to endpoint /health-check
			req := httptest.NewRequest("GET", "/health-check", nil)
			respRec := httptest.NewRecorder()
			apiEngine.ServeHTTP(respRec, req)

			// assert response
			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessage)
			
			// verify response JSON
			if tc.verifyResponse != nil {
				tc.verifyResponse(t, respRec.Body.String())
			}
		})
	}
}

package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9" 
	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck (t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode      int
		expectedMessage string
	}{
		{
			name: "Get response successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/health-check", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusOK,
			expectedMessage: "OK",
		},
		{
			name: "Wrong response",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/heath-check", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s, err := miniredis.Run()
			assert.NoError(t, err)
			defer s.Close() 
			
			redisClient := redis.NewClient(&redis.Options{
				Addr: s.Addr(),
			})
			defer redisClient.Close()

			apiEngine := api.NewEngine(&api.Config{}, redisClient)

			rec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedMessage)
		})
	}

}

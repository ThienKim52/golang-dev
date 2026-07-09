package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGenPass(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode      int
		expectedResponseContain string
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
			expectedResponseContain: `{"password":`,
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
			redisClient := redis.NewClient(&redis.Options{})
			apiEngine := api.NewEngine(&api.Config{}, redisClient)

			rec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if len(tc.expectedResponseContain) > 0 {
				assert.Equal(t, len(rec.Body.String()), len(tc.expectedResponseContain)+12+3)
			}
			assert.Contains(t, rec.Body.String(), tc.expectedResponseContain)
		})
	}

}

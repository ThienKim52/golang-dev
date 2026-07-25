package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThienKim52/golang-dev/internal/api"
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
			name: "Generate password successfully",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/genpass", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatusCode:      http.StatusOK,
			expectedResponseContain: `{"password":`,
		},
		{
			name: "Wrong endpoint",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/genpass2", nil)
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

			rec := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if len(tc.expectedResponseContain) > 0 {
				assert.Equal(t, len(rec.Body.String()), len(tc.expectedResponseContain)+12+3)
			}
			assert.Contains(t, rec.Body.String(), tc.expectedResponseContain)
		})
	}

}

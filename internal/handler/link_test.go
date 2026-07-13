package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLinkService is a mock implementation of LinkService
type MockLinkService struct {
	mock.Mock
}

func (m *MockLinkService) ShortenURL(ctx context.Context, url string, exp time.Duration) (string, error) {
	args := m.Called(ctx, url, exp)
	if args.Get(1) != nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}
func TestShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		reqBody            string
		setupMock          func(m *MockLinkService)         // callback mock case
		expectedStatusCode int
	}{
		{
			name:    "Success",
			reqBody: `{"url": "https://example.com", "exp": 604800}`,
			setupMock: func(m *MockLinkService) {
				m.On("ShortenURL", mock.Anything, "https://example.com", time.Duration(604800)).
					Return("abc1234", nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid JSON",
			reqBody:            `{"url": "invalid json`,
			setupMock:          func(m *MockLinkService) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Missing Fields",
			reqBody:            `{"url": "https://example.com"}`, 
			setupMock:          func(m *MockLinkService) {}, 
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Service Error",
			reqBody: `{"url": "https://example.com", "exp": 604800}`,
			setupMock: func(m *MockLinkService) {
				m.On("ShortenURL", mock.Anything, "https://example.com", time.Duration(604800)).
					Return("", assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init mock service for each case
			mockService := new(MockLinkService)
			tc.setupMock(mockService)

			handler := NewLinkHandler(mockService)

			// init mock request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/v1/links/shorten", bytes.NewBufferString(tc.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// call Handler
			handler.ShortenURL(c)

			// check
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

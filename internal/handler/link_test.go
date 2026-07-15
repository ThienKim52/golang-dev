package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks_linkservice "github.com/ThienKim52/golang-dev/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		reqBody            string
		setupMock          func(c *gin.Context) *mocks_linkservice.LinkService // callback mock case
		expectedStatusCode int
	}{
		{
			name:    "Success",
			reqBody: `{"url": "https://example.com", "exp": 604800}`,
			setupMock: func(c *gin.Context) *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				mockSvc.On("ShortenURL", c, "https://example.com", time.Duration(604800)*time.Second).
					Return("abc1234", nil)
				return mockSvc
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "Invalid JSON",
			reqBody: `{"url": "invalid json`,
			setupMock: func(c *gin.Context) *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				return mockSvc
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Missing Fields",
			reqBody: `{"url": "https://example.com"}`,
			setupMock: func(c *gin.Context) *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				return mockSvc
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Service Error",
			reqBody: `{"url": "https://example.com", "exp": 604800}`,
			setupMock: func(c *gin.Context) *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				mockSvc.On("ShortenURL", c, "https://example.com", time.Duration(604800)*time.Second).
					Return("", assert.AnError)
				return mockSvc
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init mock service for each case

			// init mock request
			// init mock request - Cho đoạn này lên trước, tạo ctx
			w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

			// init mock service for each case
			mockService := tc.setupMock(c)

			handler := NewLinkHandler(mockService)
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

func TestRedirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		reqBody string
		setupMock          func() *mocks_linkservice.LinkService // callback mock case
		setupTestRequest func(c *gin.Context)
		expectedStatusCode int
		expectedResponse string
	}{
		{
			name:    "Success",
			reqBody: `{"url": "https://example.com", "exp": 604800}`,
			setupMock: func() *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				mockSvc.On("ShortenURL", mock.Anything, "https://example.com", time.Duration(604800)*time.Second).
					Return("abc1234", nil)
				return mockSvc
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "Invalid JSON",
			reqBody: `{"url": "invalid json`,
			setupMock: func() *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				return mockSvc
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Missing Fields",
			reqBody: `{"url": "https://example.com"}`,
			setupMock: func() *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				return mockSvc
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Service Error",
			reqBody: `{"url": "https://example.com", "exp": 604800}`,
			setupMock: func() *mocks_linkservice.LinkService {
				mockSvc := mocks_linkservice.NewLinkService(t)
				mockSvc.On("ShortenURL", mock.Anything, "https://example.com", time.Duration(604800)*time.Second).
					Return("", assert.AnError)
				return mockSvc
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init mock service for each case

			// init mock request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/v1/links/shorten", bytes.NewBufferString(tc.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")
			mockService := tc.setupMock()
			handler := NewLinkHandler(mockService)

			// call Handler
			handler.ShortenURL(c)

			// check
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

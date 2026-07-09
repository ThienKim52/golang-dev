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

func TestNewLinkHandler(t *testing.T) {
	mockService := new(MockLinkService)
	handler := NewLinkHandler(mockService)
	assert.NotNil(t, handler)
}

func TestShortenURL_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockLinkService)
	handler := NewLinkHandler(mockService)

	reqBody := `{"url": "https://example.com", "exp": 604800}`
	expectedCode := "abc1234"

	mockService.On("ShortenURL", mock.Anything, "https://example.com", time.Duration(604800)*time.Second).Return(expectedCode, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/links/shorten", bytes.NewBufferString(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ShortenURL(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestShortenURL_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockLinkService)
	handler := NewLinkHandler(mockService)

	reqBody := `{"url": "invalid json`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/links/shorten", bytes.NewBufferString(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ShortenURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShortenURL_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockLinkService)
	handler := NewLinkHandler(mockService)

	reqBody := `{"url": "https://example.com"}` // Missing exp field

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/links/shorten", bytes.NewBufferString(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ShortenURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShortenURL_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockLinkService)
	handler := NewLinkHandler(mockService)

	reqBody := `{"url": "https://example.com", "exp": 604800}`

	mockService.On("ShortenURL", mock.Anything, "https://example.com", time.Duration(604800)*time.Second).Return("", assert.AnError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/v1/links/shorten", bytes.NewBufferString(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.ShortenURL(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

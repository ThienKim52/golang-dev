package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHealthCheckService is a mock implementation of HealthCheckService
type MockHealthCheckService struct {
	mock.Mock
}

func (m *MockHealthCheckService) GetHealthCheck(ctx context.Context, serviceName, instanceID string) (string, error) {
	args := m.Called(ctx, serviceName, instanceID)
	if args.Get(1) != nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

func TestGetResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		serviceName  string
		instanceID   string
		mockResponse string
		mockError    error
		expectedCode int
	}{
		{
			name:         "Successful health check",
			serviceName:  "test-service",
			instanceID:   "123e4567-e89b-12d3-a456-426614174000",
			mockResponse: "OK",
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Health check with empty values",
			serviceName:  "",
			instanceID:   "",
			mockResponse: "OK",
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Redis connection error",
			serviceName:  "test-service",
			instanceID:   "123e4567-e89b-12d3-a456-426614174000",
			mockResponse: "",
			mockError:    assert.AnError,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockHealthCheckService)
			mockService.On("GetHealthCheck", mock.Anything, tt.serviceName, tt.instanceID).Return(tt.mockResponse, tt.mockError)

			// Create handler
			handler := NewHealthCheck(mockService, tt.serviceName, tt.instanceID)

			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/health-check", nil)
			c.Set("service_name", tt.serviceName)
			c.Set("instance_id", tt.instanceID)

			// Call handler
			handler.GetResponse(c)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

// Mock data test heathcheck
func TestNewHealthCheck(t *testing.T) {
	mockService := new(MockHealthCheckService)
	handler := NewHealthCheck(mockService, "test-service", "123e4567-e89b-12d3-a456-426614174000")
	assert.NotNil(t, handler)
}

package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHealthCheckService is a mock implementation of HealthCheckService
type MockHealthCheckService struct {
	mock.Mock
}
// Test healthcheck services
func (m *MockHealthCheckService) GetHealthCheck(ctx context.Context) (service.HealthCheckResult, error) {
	args := m.Called(ctx)
	return args.Get(0).(service.HealthCheckResult), args.Error(1)
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
			result := service.HealthCheckResult{
				Message: tt.mockResponse,
				ServiceName: tt.serviceName,
				InstanceID:  tt.instanceID,
			}
			mockService.On("GetHealthCheck", mock.Anything).Return(result, tt.mockError)

			// Create handler
			handler := NewHealthCheck(mockService)

			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/health-check", nil)

			// Call handler
			handler.GetResponse(c)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

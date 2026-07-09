package service

import (
	"context"
	"testing"

	"github.com/ThienKim52/golang-dev/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthCheck(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		serviceName string
		instanceID  string
		expected    string
		mockPingErr error
		expectError bool
	}{
		{
			name:        "Valid service name and instance ID",
			serviceName: "test-service",
			instanceID:  "123e4567-e89b-12d3-a456-426614174000",
			expected:    "OK",
			mockPingErr: nil,
			expectError: false,
		},
		{
			name:        "Empty service name and instance ID",
			serviceName: "",
			instanceID:  "",
			expected:    "OK",
			mockPingErr: nil,
			expectError: false,
		},
		{
			name:        "Redis ping error",
			serviceName: "test-service",
			instanceID:  "123e4567-e89b-12d3-a456-426614174000",
			expected:    "",
			mockPingErr: assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockLinkRepository)
			service := NewHealthCheck(mockRepo)
			mockRepo.On("Ping", ctx).Return(tt.mockPingErr)

			result, err := service.GetHealthCheck(ctx, tt.serviceName, tt.instanceID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewHealthCheck(t *testing.T) {
	mockRepo := new(repository.MockLinkRepository)
	service := NewHealthCheck(mockRepo)
	assert.NotNil(t, service)
	_, ok := service.(*healthCheckService)
	assert.True(t, ok)
}

func TestNewHealthCheckWithoutRedis(t *testing.T) {
	service := NewHealthCheckWithoutRedis()
	assert.NotNil(t, service)

	ctx := context.Background()
	result, err := service.GetHealthCheck(ctx, "test-service", "test-id")

	assert.NoError(t, err)
	assert.Equal(t, "OK", result)
}

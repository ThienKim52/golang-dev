package service

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/repository"
)

// HealthCheckService defines the interface for health check operations
type HealthCheckService interface {
	GetHealthCheck(ctx context.Context, serviceName, instanceID string) (message string, err error)
}

// healthCheckService implements HealthCheckService
type healthCheckService struct {
	repo repository.LinkRepository
}

// NewHealthCheck creates a new health check service
func NewHealthCheck(repo repository.LinkRepository) HealthCheckService {
	return &healthCheckService{
		repo: repo,
	}
}

// NewHealthCheckWithoutRedis creates a new health check service without Redis dependency
func NewHealthCheckWithoutRedis() HealthCheckService {
	return &healthCheckService{
		repo: nil,
	}
}

// GetHealthCheck returns the health check response
func (s *healthCheckService) GetHealthCheck(ctx context.Context, serviceName, instanceID string) (string, error) {
	// Check Redis connection if repository is available
	if s.repo != nil {
		if err := s.repo.Ping(ctx); err != nil {
			return "", err
		}
	}

	return "OK", nil
}

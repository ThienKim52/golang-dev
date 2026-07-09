package service

import (
	"context"
	"testing"
	"time"

	"github.com/ThienKim52/golang-dev/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGenPass is a mock implementation of GenPass
type MockGenPass struct {
	mock.Mock
}

func (m *MockGenPass) GeneratePassword(passwordLength int) (string, error) {
	args := m.Called(passwordLength)
	return args.String(0), args.Error(1)
}

func TestNewLinkService(t *testing.T) {
	mockRepo := new(repository.MockLinkRepository)
	mockGenPass := new(MockGenPass)
	service := NewLinkService(mockRepo, mockGenPass)
	assert.NotNil(t, service)
	_, ok := service.(*linkService)
	assert.True(t, ok)
}

func TestShortenURL_Success(t *testing.T) {
	mockRepo := new(repository.MockLinkRepository)
	mockGenPass := new(MockGenPass)
	service := NewLinkService(mockRepo, mockGenPass)
	ctx := context.Background()

	url := "https://example.com"
	exp := time.Duration(604800) * time.Second

	// Mock genPass to return a code
	mockGenPass.On("GeneratePassword", codeLength).Return("abc1234", nil)
	// Mock repository to return false for exists check
	mockRepo.On("Exists", ctx, "abc1234").Return(false, nil)
	// Mock repository to save successfully
	mockRepo.On("Save", ctx, "abc1234", url, exp).Return(nil)

	code, err := service.ShortenURL(ctx, url, exp)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)
	assert.Equal(t, "abc1234", code)
	assert.Len(t, code, codeLength)

	mockRepo.AssertExpectations(t)
	mockGenPass.AssertExpectations(t)
}

func TestShortenURL_CodeConflict_RetrySuccess(t *testing.T) {
	mockRepo := new(repository.MockLinkRepository)
	mockGenPass := new(MockGenPass)
	service := NewLinkService(mockRepo, mockGenPass)
	ctx := context.Background()

	url := "https://example.com"
	exp := time.Duration(604800) * time.Second

	// First code exists, second code doesn't
	mockGenPass.On("GeneratePassword", codeLength).Return("code1", nil).Once()
	mockRepo.On("Exists", ctx, "code1").Return(true, nil).Once()
	mockGenPass.On("GeneratePassword", codeLength).Return("code2", nil).Once()
	mockRepo.On("Exists", ctx, "code2").Return(false, nil).Once()
	mockRepo.On("Save", ctx, "code2", url, exp).Return(nil)

	code, err := service.ShortenURL(ctx, url, exp)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)
	assert.Equal(t, "code2", code)

	mockRepo.AssertExpectations(t)
	mockGenPass.AssertExpectations(t)
}

func TestShortenURL_MaxRetriesExceeded(t *testing.T) {
	mockRepo := new(repository.MockLinkRepository)
	mockGenPass := new(MockGenPass)
	service := NewLinkService(mockRepo, mockGenPass)
	ctx := context.Background()

	url := "https://example.com"
	exp := time.Duration(604800) * time.Second

	// All codes exist
	mockGenPass.On("GeneratePassword", codeLength).Return("code", nil)
	mockRepo.On("Exists", ctx, mock.AnythingOfType("string")).Return(true, nil)

	code, err := service.ShortenURL(ctx, url, exp)

	assert.Error(t, err)
	assert.Empty(t, code)
	assert.IsType(t, &ErrMaxRetriesExceeded{}, err)

	mockRepo.AssertExpectations(t)
	mockGenPass.AssertExpectations(t)
}

func TestShortenURL_RepositoryError(t *testing.T) {
	mockRepo := new(repository.MockLinkRepository)
	mockGenPass := new(MockGenPass)
	service := NewLinkService(mockRepo, mockGenPass)
	ctx := context.Background()

	url := "https://example.com"
	exp := time.Duration(604800) * time.Second

	// Mock genPass to return a code
	mockGenPass.On("GeneratePassword", codeLength).Return("abc1234", nil)
	// Repository returns error on exists check
	mockRepo.On("Exists", ctx, "abc1234").Return(false, assert.AnError)

	code, err := service.ShortenURL(ctx, url, exp)

	assert.Error(t, err)
	assert.Empty(t, code)

	mockRepo.AssertExpectations(t)
	mockGenPass.AssertExpectations(t)
}

package service

import (
	"context"
	"testing"
	"time"
	"errors"

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

func TestLinkService_ShortenURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		exp          time.Duration
		setupMocks   func(mockRepo *repository.MockLinkRepository, mockGenPass *MockGenPass)
		expectedCode string
		expectedErr  error
	}{
		{
			name: "Success",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *repository.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("abc1234", nil)
				mockRepo.On("Exists", ctx, "abc1234").Return(false, nil)
				mockRepo.On("Save", ctx, "abc1234", "https://example.com", 604800*time.Second).Return(nil)
			},
			expectedCode: "abc1234",
			expectedErr:  nil,
		},
		{
			name: "CodeConflict_RetrySuccess",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *repository.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("code1", nil).Once()
				mockRepo.On("Exists", ctx, "code1").Return(true, nil).Once()
				mockGenPass.On("GeneratePassword", codeLength).Return("code2", nil).Once()
				mockRepo.On("Exists", ctx, "code2").Return(false, nil).Once()
				mockRepo.On("Save", ctx, "code2", "https://example.com", 604800*time.Second).Return(nil)
			},
			expectedCode: "code2",
			expectedErr:  nil,
		},
		{
			name: "MaxRetriesExceeded",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *repository.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("code", nil)
				mockRepo.On("Exists", ctx, mock.AnythingOfType("string")).Return(true, nil)
			},
			expectedCode: "",
			expectedErr:  ErrMaxRetriesExceeded,
		},
		{
			name: "RepositoryError",
			url:  "https://example.com",
			exp:  604800 * time.Second,
			setupMocks: func(mockRepo *repository.MockLinkRepository, mockGenPass *MockGenPass) {
				ctx := context.Background()
				mockGenPass.On("GeneratePassword", codeLength).Return("abc1234", nil)
				mockRepo.On("Exists", ctx, "abc1234").Return(false, errors.New("repository error"))
			},
			expectedCode: "",
			expectedErr:  errors.New("repository error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(repository.MockLinkRepository)
			mockGenPass := new(MockGenPass)
			service := NewLinkService(mockRepo, mockGenPass)

			tc.setupMocks(mockRepo, mockGenPass)

			ctx := context.Background()
			code, err := service.ShortenURL(ctx, tc.url, tc.exp)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, code)
			}

			mockRepo.AssertExpectations(t)
			mockGenPass.AssertExpectations(t)
		})
	}
}
	

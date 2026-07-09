package repository

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


// MockRedisClient is a mock implementation of the private redisClient interface
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		val := int64(0)
		if args.Get(0) != nil {
			val = args.Get(0).(int64)
		}
		cmd.SetVal(val)
	}
	return cmd
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	}
	return cmd
}

func TestNewRedisLinkRepository(t *testing.T) {
	mockClient := new(MockRedisClient)
	repo := NewRedisLinkRepository(mockClient)
	assert.NotNil(t, repo)
	_, ok := repo.(*RedisLinkRepository)
	assert.True(t, ok)
}

func TestRedisLinkRepository_Save(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		url         string
		exp         time.Duration
		setupMock   func(m *MockRedisClient)
		expectedErr error
	}{
		{
			name: "save successful",
			code: "test",
			url:  "http://google.com",
			exp:  0,
			setupMock: func(m *MockRedisClient) {
				m.On("Set", mock.Anything, "test", "http://google.com", time.Duration(0)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "save error",
			code: "test",
			url:  "http://google.com",
			exp:  0,
			setupMock: func(m *MockRedisClient) {
				m.On("Set", mock.Anything, "test", "http://google.com", time.Duration(0)).Return(assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockClient := new(MockRedisClient)
			tc.setupMock(mockClient)

			repo := NewRedisLinkRepository(mockClient)
			err := repo.Save(t.Context(), tc.code, tc.url, tc.exp)
			assert.Equal(t, tc.expectedErr, err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRedisLinkRepository_GetByCode(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		code        string
		setupMock   func(m *MockRedisClient)
		expectedVal string
		expectedErr error
	}{
		{
			name: "get successful",
			code: "test",
			setupMock: func(m *MockRedisClient) {
				m.On("Get", mock.Anything, "test").Return("http://google.com", nil)
			},
			expectedVal: "http://google.com",
			expectedErr: nil,
		},
		{
			name: "get error/not found",
			code: "test",
			setupMock: func(m *MockRedisClient) {
				m.On("Get", mock.Anything, "test").Return("", redis.Nil)
			},
			expectedVal: "",
			expectedErr: redis.Nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockClient := new(MockRedisClient)
			tc.setupMock(mockClient)

			repo := NewRedisLinkRepository(mockClient)
			val, err := repo.GetByCode(t.Context(), tc.code)
			assert.Equal(t, tc.expectedVal, val)
			assert.Equal(t, tc.expectedErr, err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRedisLinkRepository_Exists(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name           string
		code           string
		setupMock      func(m *MockRedisClient)
		expectedExists bool
		expectedErr    error
	}{
		{
			name: "exists true",
			code: "test",
			setupMock: func(m *MockRedisClient) {
				m.On("Exists", mock.Anything, []string{"test"}).Return(int64(1), nil)
			},
			expectedExists: true,
			expectedErr:    nil,
		},
		{
			name: "exists false",
			code: "test",
			setupMock: func(m *MockRedisClient) {
				m.On("Exists", mock.Anything, []string{"test"}).Return(int64(0), nil)
			},
			expectedExists: false,
			expectedErr:    nil,
		},
		{
			name: "exists error",
			code: "test",
			setupMock: func(m *MockRedisClient) {
				m.On("Exists", mock.Anything, []string{"test"}).Return(int64(0), assert.AnError)
			},
			expectedExists: false,
			expectedErr:    assert.AnError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockClient := new(MockRedisClient)
			tc.setupMock(mockClient)

			repo := NewRedisLinkRepository(mockClient)
			exists, err := repo.Exists(t.Context(), tc.code)
			assert.Equal(t, tc.expectedExists, exists)
			assert.Equal(t, tc.expectedErr, err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRedisLinkRepository_Ping(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name        string
		setupMock   func(m *MockRedisClient)
		expectedErr error
	}{
		{
			name: "ping successful",
			setupMock: func(m *MockRedisClient) {
				m.On("Ping", mock.Anything).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "ping error",
			setupMock: func(m *MockRedisClient) {
				m.On("Ping", mock.Anything).Return(assert.AnError)
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockClient := new(MockRedisClient)
			tc.setupMock(mockClient)

			repo := NewRedisLinkRepository(mockClient)
			err := repo.Ping(t.Context())
			assert.Equal(t, tc.expectedErr, err)
			mockClient.AssertExpectations(t)
		})
	}
}

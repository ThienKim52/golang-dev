package service

import (
	"context"
	"time"

	"github.com/ThienKim52/golang-dev/internal/repository"
)

const (
	codeLength = 7
	maxRetries = 10
)

// LinkService defines the interface for link operations
type LinkService interface {
	ShortenURL(ctx context.Context, url string, exp time.Duration) (code string, err error)
}

// linkService implements LinkService
type linkService struct {
	repo    repository.LinkRepository
	genPass GenPass
}

// NewLinkService creates a new link service
func NewLinkService(repo repository.LinkRepository, genPass GenPass) LinkService {
	return &linkService{
		repo:    repo,
		genPass: genPass,
	}
}

// ShortenURL creates a short code for a URL and saves it to the repository
func (s *linkService) ShortenURL(ctx context.Context, url string, exp time.Duration) (string, error) {
	var code string
	var err error

	// Generate a unique code with retry mechanism
	for i := 0; i < maxRetries; i++ {
		code, err = s.genPass.GeneratePassword(codeLength)
		if err != nil {
			return "", err
		}

		// Check if code already exists
		exists, err := s.repo.Exists(ctx, code)
		if err != nil {
			return "", err
		}

		if !exists {
			// Code is unique, save the link
			if err := s.repo.Save(ctx, code, url, exp); err != nil {
				return "", err
			}

			return code, nil
		}
		// Code exists, try again
	}

	return "", &ErrMaxRetriesExceeded{maxRetries: maxRetries}
}

// ErrMaxRetriesExceeded is returned when max retries for generating unique code is exceeded
type ErrMaxRetriesExceeded struct {
	maxRetries int
}

func (e *ErrMaxRetriesExceeded) Error() string {
	return "max retries exceeded while generating unique code"
}

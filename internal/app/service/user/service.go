package user

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/model"
	"github.com/ThienKim52/golang-dev/internal/app/repository/user"
	"github.com/ThienKim52/golang-dev/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

type service struct {
	repo user.Repository
	hasher utils.Hasher
}

func NewService(repo user.Repository, hasher utils.Hasher)Service{
	return &service{
		repo: repo,
		hasher: hasher,
	}
}
package user
import (
	"context"
	"github.com/ThienKim52/golang-dev/internal/app/model"
)
type Repository interface{
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
}
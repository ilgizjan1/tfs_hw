package domain

import (
	"context"
	"homework10/internal/domain/models"
)

//go:generate mockgen -source=./user.go -destination=../service/mock/user.go -package=repoMock UserRepository
type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	AddUser(ctx context.Context, user models.User) (int64, error)
	Update(ctx context.Context, userID int64, nickName string, email string) (*models.User, error)
	Delete(ctx context.Context, userID int64) error
}

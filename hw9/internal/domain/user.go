package domain

import (
	"context"
	"homework9/internal/domain/models"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	AddUser(ctx context.Context, user models.User) (int64, error)
	Update(ctx context.Context, userID int64, nickName string, email string) (*models.User, error)
	Delete(ctx context.Context, userID int64) error
}

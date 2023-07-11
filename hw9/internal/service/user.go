package service

import (
	"context"
	"homework9/internal/domain"
	"homework9/internal/domain/models"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, nickName string, email string) (*models.User, error) {
	user := models.User{NickName: nickName, Email: email}
	userID, err := s.userRepo.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID int64, nickName string, email string) (*models.User, error) {
	user, err := s.userRepo.Update(ctx, userID, nickName, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}

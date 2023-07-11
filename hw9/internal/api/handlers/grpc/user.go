package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	contracts "homework9/internal/api/handlers/grpc/contracts/langs/go"
	"homework9/internal/api/handlers/grpc/mapper"
	"homework9/internal/domain/models"
)

type UserService interface {
	CreateUser(ctx context.Context, nickName string, email string) (*models.User, error)
	UpdateUser(ctx context.Context, userID int64, nickName string, email string) (*models.User, error)
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	DeleteUser(ctx context.Context, userID int64) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userServ UserService) *UserHandler {
	return &UserHandler{
		userService: userServ,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, request *contracts.CreateUserRequest) (*contracts.UserResponse, error) {
	user, err := h.userService.CreateUser(ctx, request.Nickname, request.Email)
	if err != nil {
		return nil, err
	}
	return mapper.UserToResponse(user), nil
}

func (h *UserHandler) GetUser(ctx context.Context, request *contracts.GetUserRequest) (*contracts.UserResponse, error) {
	user, err := h.userService.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return mapper.UserToResponse(user), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, request *contracts.UpdateUserRequest) (*contracts.UserResponse, error) {
	user, err := h.userService.UpdateUser(ctx, request.UserId, request.Nickname, request.Email)
	if err != nil {
		return nil, err
	}
	return mapper.UserToResponse(user), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, request *contracts.DeleteUserRequest) (*emptypb.Empty, error) {
	err := h.userService.DeleteUser(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

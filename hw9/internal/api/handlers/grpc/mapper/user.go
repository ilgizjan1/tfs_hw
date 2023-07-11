package mapper

import (
	contracts "homework9/internal/api/handlers/grpc/contracts/langs/go"
	"homework9/internal/domain/models"
)

func UserToResponse(user *models.User) *contracts.UserResponse {
	return &contracts.UserResponse{
		UserId:   user.ID,
		Nickname: user.NickName,
		Email:    user.Email,
	}
}

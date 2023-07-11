package mapper

import (
	contracts "homework10/internal/api/handlers/grpc/contracts/langs/go"
	"homework10/internal/domain/models"
)

func UserToResponse(user *models.User) *contracts.UserResponse {
	return &contracts.UserResponse{
		UserId:   user.ID,
		Nickname: user.NickName,
		Email:    user.Email,
	}
}

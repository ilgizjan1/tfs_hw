package mapper

import (
	"homework10/internal/api/handlers/httpgin/response"
	"homework10/internal/domain/models"

	"github.com/gofiber/fiber/v2"
)

func UserSuccessResponse(user *models.User) *fiber.Map {
	return &fiber.Map{
		"data": response.UserResponse{
			ID:       user.ID,
			Nickname: user.NickName,
			Email:    user.Email,
		},
	}
}

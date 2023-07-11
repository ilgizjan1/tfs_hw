package mapper

import (
	"github.com/gofiber/fiber/v2"
	"homework9/internal/api/handlers/httpgin/response"
	"homework9/internal/domain/models"
)

func UserSuccessResponse(user *models.User) *fiber.Map {
	return &fiber.Map{
		"data": response.UserResponse{
			ID:       user.ID,
			Nickname: user.NickName,
			Email:    user.Email,
		},
		"error": nil,
	}
}

package mapper

import (
	"testing"

	"homework10/internal/api/handlers/httpgin/response"
	"homework10/internal/domain/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestUserSuccessResponse(t *testing.T) {

	tests := []struct {
		name     string
		user     *models.User
		expected *fiber.Map
	}{
		{
			name: "successfully map user to response",
			user: &models.User{
				ID:       0,
				NickName: "test_user",
				Email:    "test_user@gmail.com",
			},
			expected: &fiber.Map{
				"data": response.UserResponse{
					ID:       0,
					Nickname: "test_user",
					Email:    "test_user@gmail.com",
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			actual := UserSuccessResponse(tc.user)

			require.Equal(t, tc.expected, actual)
		})
	}
}

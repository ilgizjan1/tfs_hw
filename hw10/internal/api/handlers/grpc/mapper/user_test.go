package mapper

import (
	contracts "homework10/internal/api/handlers/grpc/contracts/langs/go"
	"homework10/internal/domain/models"
	"reflect"
	"testing"
)

func TestUserToResponse(t *testing.T) {
	tests := []struct {
		name     string
		user     *models.User
		expected *contracts.UserResponse
	}{
		{
			name: "successfully map user to response",
			user: &models.User{
				ID:       0,
				NickName: "test_user",
				Email:    "test_user@gmail.com",
			},
			expected: &contracts.UserResponse{
				UserId:   0,
				Nickname: "test_user",
				Email:    "test_user@gmail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserToResponse(tt.user); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("UserToResponse() = %v, want %v", got, tt.expected)
			}
		})
	}
}

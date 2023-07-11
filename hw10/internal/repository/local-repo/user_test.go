package localrepo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"homework10/internal/domain/models"
	"testing"
)

func TestUserRepo_GetUser(t *testing.T) {
	userRepo := NewUserRepo()
	userRepo.storage[0] = &models.User{
		ID:       0,
		NickName: "test nickname",
		Email:    "test email",
	}
	tests := []struct {
		name         string
		userID       int64
		expectedUser *models.User
		err          error
		cancel       bool
	}{
		{
			name:   "successfully test GetUser()",
			userID: 0,
			expectedUser: &models.User{
				ID:       0,
				NickName: "test nickname",
				Email:    "test email",
			},
			err:    nil,
			cancel: false,
		},
		{
			name:   " test GetUser()",
			userID: 10,
			err:    fmt.Errorf("the user does not exist"),
			cancel: false,
		},
		{
			name:   "cancel context test GetUser()",
			userID: 0,
			expectedUser: &models.User{
				ID:       0,
				NickName: "test nickname",
				Email:    "test email",
			},
			err:    context.Canceled,
			cancel: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			user, err := userRepo.GetUser(ctx, tc.userID)
			if err != nil {
				assert.Equal(t, tc.err, err)
				assert.Nil(t, user)
			} else {
				assert.Equal(t, *tc.expectedUser, *user)
			}
		})
	}
}

func TestUserRepo_AddUser(t *testing.T) {
	userRepo := NewUserRepo()

	tests := []struct {
		name       string
		newUser    models.User
		expectedID int64
		err        error
		cancel     bool
	}{
		{
			name: "successfully test GetUser()",
			newUser: models.User{
				ID:       0,
				NickName: "test nickname",
				Email:    "test email",
			},
			expectedID: 0,
			err:        nil,
			cancel:     false,
		},
		{
			name: "cancel context test GetUser()",
			newUser: models.User{
				ID:       0,
				NickName: "test nickname",
				Email:    "test email",
			},
			err:        context.Canceled,
			expectedID: 0,
			cancel:     true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			userID, err := userRepo.AddUser(ctx, tc.newUser)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedID, userID)
		})
	}
}

func TestUserRepo_UpdateUser(t *testing.T) {
	userRepo := NewUserRepo()
	userRepo.storage[0] = &models.User{
		ID:       0,
		NickName: "test nickname",
		Email:    "test email",
	}
	tests := []struct {
		name         string
		expectedUser *models.User
		err          error
		cancel       bool
	}{
		{
			name: "successfully test GetUser()",
			expectedUser: &models.User{
				ID:       0,
				NickName: "new nickname",
				Email:    "new email",
			},
			err:    nil,
			cancel: false,
		},
		{
			name: "cancel context test GetUser()",
			expectedUser: &models.User{
				ID:       0,
				NickName: "new nickname",
				Email:    "new email",
			},
			err:    context.Canceled,
			cancel: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			user, err := userRepo.Update(ctx, 0, tc.expectedUser.NickName, tc.expectedUser.Email)
			if err != nil {
				assert.Equal(t, tc.err, err)
				assert.Nil(t, user)
			} else {
				assert.Equal(t, *tc.expectedUser, *user)
			}
		})
	}
}

func TestUserRepo_DeleteUser(t *testing.T) {
	userRepo := NewUserRepo()
	userRepo.storage[0] = &models.User{
		ID:       0,
		NickName: "test nickname",
		Email:    "test email",
	}
	tests := []struct {
		name   string
		err    error
		cancel bool
	}{
		{
			name:   "successfully test GetUser()",
			err:    nil,
			cancel: false,
		},
		{
			name:   "cancel context test GetUser()",
			err:    context.Canceled,
			cancel: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			err := userRepo.Delete(ctx, 0)
			if err != nil {
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, len(userRepo.storage), 0)
			}
		})
	}
}

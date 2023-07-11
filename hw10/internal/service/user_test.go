package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"homework10/internal/domain/models"
	repoMock "homework10/internal/service/mock"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := repoMock.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo)

	testTable := []struct {
		name    string
		userID  int64
		inUser  *models.User
		outUser *models.User
		repoErr error
		wantErr bool
	}{
		{
			name:   "true test CreateUser()",
			userID: 100,
			inUser: &models.User{
				NickName: "Ivan",
				Email:    "Ivan@gmail.com",
			},
			outUser: &models.User{
				ID:       100,
				NickName: "Ivan",
				Email:    "Ivan@gmail.com",
			},
			wantErr: false,
		},
		{
			name:   "error from repository GetUser()",
			userID: 0,
			inUser: &models.User{
				NickName: "Ivan",
				Email:    "Ivan@gmail.com",
			},
			repoErr: fmt.Errorf("error from repo GetUser()"),
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo.EXPECT().AddUser(ctx, *testCase.inUser).Return(testCase.userID, testCase.repoErr).Times(1)

			user, err := userService.CreateUser(ctx, testCase.inUser.NickName, testCase.inUser.Email)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.outUser, *user)
		})
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := repoMock.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo)

	testTable := []struct {
		name     string
		userID   int64
		expected *models.User
		repoErr  error
		wantErr  bool
	}{
		{
			name:   "true test",
			userID: 100,
			expected: &models.User{
				ID:       100,
				NickName: "Ivan",
				Email:    "Ivan@gmail.com",
			},
			wantErr: false,
		},
		{
			name:     "error from repo GetUser()",
			userID:   47,
			expected: nil,
			repoErr:  fmt.Errorf("error from GetUser()"),
			wantErr:  true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo.EXPECT().GetUser(ctx, testCase.userID).Return(testCase.expected, testCase.repoErr).Times(1)

			user, err := userService.GetUser(ctx, testCase.userID)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.expected, *user)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := repoMock.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo)

	testTable := []struct {
		name     string
		inUser   *models.User
		expected *models.User
		repoErr  error
		wantErr  bool
	}{
		{
			name: "true test",
			inUser: &models.User{
				ID:       100,
				NickName: "ivan",
				Email:    "ivan@gmail.com",
			},
			expected: &models.User{
				ID:       100,
				NickName: "ivan",
				Email:    "ivan@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "error from GetUser()",
			inUser: &models.User{
				ID:       100,
				NickName: "ivan",
				Email:    "ivan@gmail.com",
			},
			expected: nil,
			repoErr:  fmt.Errorf("error from GetUser()"),
			wantErr:  true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo.EXPECT().Update(ctx, testCase.inUser.ID, testCase.inUser.NickName, testCase.inUser.Email).
				Return(testCase.expected, testCase.repoErr).Times(1)

			user, err := userService.UpdateUser(
				ctx,
				testCase.inUser.ID,
				testCase.inUser.NickName,
				testCase.inUser.Email,
			)

			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *testCase.expected, *user)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := repoMock.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo)

	testTable := []struct {
		name      string
		userID    int64
		getAdErr  error
		deleteErr error
		wantErr   bool
	}{
		{
			name:    "true test",
			userID:  100,
			wantErr: false,
		},
		{
			name:     "error from GetUser()",
			userID:   100,
			getAdErr: fmt.Errorf("error from GetUser()"),
			wantErr:  true,
		},
		{
			name:      "error from DeleteUser()",
			userID:    100,
			deleteErr: fmt.Errorf("error from DeleteUser()"),
			wantErr:   true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo.EXPECT().GetUser(ctx, testCase.userID).
				Return(nil, testCase.getAdErr).Times(1)

			if testCase.getAdErr == nil {
				userRepo.EXPECT().Delete(ctx, testCase.userID).Return(testCase.deleteErr).Times(1)
			}

			err := userService.DeleteUser(ctx, testCase.userID)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

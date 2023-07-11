package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homework10/internal/api/handlers/httpgin/request"
	"net/http"
	"net/http/httptest"
	"testing"

	"homework10/internal/api/handlers/httpgin/mock"
	"homework10/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_getUser(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		mockBehaviour      func(service *handlerMock.MockUserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "successfully get user",
			userID: "1",
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().
					GetUser(gomock.Any(), int64(1)).
					Return(
						&models.User{
							ID:       1,
							NickName: "test nickname",
							Email:    "test email",
						}, nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 1,
						"nickname": "test nickname",
						"email": "test email"
					}
				}
				`,
		},
		{
			name:               "invalid user id passed",
			userID:             "invalid_user_id_1",
			mockBehaviour:      func(service *handlerMock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_user_id_1\": invalid syntax"}`,
		},
		{
			name:   "error from service",
			userID: "1",
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().GetUser(gomock.Any(), int64(1)).Return(nil, fmt.Errorf("error from service"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error": "error from service"}`,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := handlerMock.NewMockUserService(ctrl)
			tc.mockBehaviour(service)

			handler := NewUserHandler(service)

			//Test Server
			rg := gin.New()
			rg.GET("/:user_id", handler.getUser)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.userID), nil)

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_createUser(t *testing.T) {
	tests := []struct {
		name               string
		user               any
		mockBehaviour      func(service *handlerMock.MockUserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully create user",
			user: request.CreateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().CreateUser(gomock.Any(), "test nickname", "test email").
					Return(&models.User{
						ID:       int64(0),
						NickName: "test nickname",
						Email:    "test email",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 0,
						"nickname": "test nickname",
						"email": "test email"
					}
				}
				`,
		},
		{
			name: "invalid json passed",
			user: struct {
				Nickname bool `json:"nickname"`
			}{
				Nickname: true,
			},
			mockBehaviour:      func(service *handlerMock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "json: cannot unmarshal bool into Go struct field CreateUserRequest.nickname of type string"}`,
		},
		{
			name: "error from service",
			user: request.CreateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().CreateUser(gomock.Any(), "test nickname", "test email").
					Return(nil, fmt.Errorf("error from service"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error": "error from service"}`,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := handlerMock.NewMockUserService(ctrl)
			tc.mockBehaviour(service)

			handler := NewUserHandler(service)

			//Test Server
			rg := gin.New()
			rg.POST("/users", handler.createUser)

			jsonValue, err := json.Marshal(tc.user)
			require.Equal(t, err, nil)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonValue))

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_updateUser(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		user               any
		mockBehaviour      func(service *handlerMock.MockUserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully update user",
			user: request.UpdateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			userID: "0",
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().UpdateUser(gomock.Any(), int64(0), "test nickname", "test email").
					Return(&models.User{
						ID:       int64(0),
						NickName: "test nickname",
						Email:    "test email",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 0,
						"nickname": "test nickname",
						"email": "test email"
					}
				}
				`,
		},
		{
			name:   "invalid json passed",
			userID: "0",
			user: struct {
				Nickname bool `json:"nickname"`
			}{
				Nickname: true,
			},
			mockBehaviour:      func(service *handlerMock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"json: cannot unmarshal bool into Go struct field UpdateUserRequest.nickname of type string"}`,
		},
		{
			name:   "error param parsing",
			userID: "invalid_param",
			user: request.CreateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			mockBehaviour:      func(service *handlerMock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_param\": invalid syntax"}`,
		},
		{
			name:   "error from service",
			userID: "0",
			user: request.CreateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().UpdateUser(gomock.Any(), int64(0), "test nickname", "test email").
					Return(nil, fmt.Errorf("error from service"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error": "error from service"}`,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := handlerMock.NewMockUserService(ctrl)
			tc.mockBehaviour(service)

			handler := NewUserHandler(service)

			//Test Server
			rg := gin.New()
			rg.PUT("/:user_id", handler.updateUser)

			jsonValue, err := json.Marshal(tc.user)
			require.Equal(t, err, nil)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s", tc.userID), bytes.NewBuffer(jsonValue))

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_deleteUser(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		user               any
		mockBehaviour      func(service *handlerMock.MockUserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "successfully update user",
			userID: "0",

			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().DeleteUser(gomock.Any(), int64(0)).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"success":"User #0 deleted"}`,
		},
		{
			name:   "error param parsing",
			userID: "invalid_param",
			user: request.CreateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			mockBehaviour:      func(service *handlerMock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_param\": invalid syntax"}`,
		},
		{
			name:   "error from service",
			userID: "0",
			user: request.CreateUserRequest{
				NickName: "test nickname",
				Email:    "test email",
			},
			mockBehaviour: func(service *handlerMock.MockUserService) {
				service.EXPECT().DeleteUser(gomock.Any(), int64(0)).
					Return(fmt.Errorf("error from service"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error": "error from service"}`,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := handlerMock.NewMockUserService(ctrl)
			tc.mockBehaviour(service)

			handler := NewUserHandler(service)

			//Test Server
			rg := gin.New()
			rg.DELETE("/:user_id", handler.deleteUser)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tc.userID), nil)

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

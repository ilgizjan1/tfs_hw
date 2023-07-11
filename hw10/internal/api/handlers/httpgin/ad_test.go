package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ilgizjan1/publication"
	"github.com/stretchr/testify/require"
	handlerMock "homework10/internal/api/handlers/httpgin/mock"
	"homework10/internal/api/handlers/httpgin/request"
	"homework10/internal/domain/models"
	"homework10/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_getAd(t *testing.T) {
	tests := []struct {
		name               string
		adID               string
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			adID: "0",
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					GetAdByID(gomock.Any(), int64(0)).
					Return(
						&models.Ad{
							ID:        0,
							Title:     "test title",
							Text:      "test text",
							UserID:    0,
							Published: true,
						}, nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 0,
						"title": "test title",
						"text": "test text",
						"user_id": 0,
						"published": true,
						"date_creation": "",
						"date_update": ""
					}
				}
				`,
		},
		{
			name:               "invalid user id passed",
			adID:               "invalid_user_id_1",
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_user_id_1\": invalid syntax"}`,
		},
		{
			name: "error from service",
			adID: "0",
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().GetAdByID(gomock.Any(), int64(0)).Return(nil, fmt.Errorf("error from service"))
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.GET("/:ad_id", handler.getAd)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.adID), nil)

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_createAd(t *testing.T) {
	tests := []struct {
		name               string
		adID               string
		ad                 any
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			ad: request.CreateAdRequest{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					CreateAd(gomock.Any(), "test title", "test text", int64(0)).
					Return(
						&models.Ad{
							ID:        0,
							Title:     "test title",
							Text:      "test text",
							UserID:    0,
							Published: true,
						}, nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 0,
						"title": "test title",
						"text": "test text",
						"user_id": 0,
						"published": true,
						"date_creation": "",
						"date_update": ""
					}
				}
				`,
		},
		{
			name: "invalid user id passed",
			ad: struct {
				Title bool `json:"title"`
			}{
				Title: true,
			},
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "json: cannot unmarshal bool into Go struct field CreateAdRequest.title of type string"}`,
		},
		{
			name: "error from service: Validation",
			ad: request.CreateAdRequest{
				Title:  "",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().CreateAd(gomock.Any(), "", "test text", int64(0)).
					Return(nil, publication.ValidationErrors{publication.ValidationError{Err: publication.ErrInvalidTitle}})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "wrong title"}`,
		},
		{
			name: "error from service",
			ad: request.CreateAdRequest{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().CreateAd(gomock.Any(), "test title", "test text", int64(0)).
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.POST("/ads", handler.createAd)

			jsonValue, err := json.Marshal(tc.ad)
			require.Equal(t, err, nil)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/ads", bytes.NewBuffer(jsonValue))

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_changeAdStatus(t *testing.T) {
	tests := []struct {
		name               string
		adID               string
		ad                 any
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			adID: "0",
			ad: request.ChangeAdStatusRequest{
				UserID:    0,
				Published: true,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					ChangeAdStatus(gomock.Any(), int64(0), int64(0), true).
					Return(
						&models.Ad{
							ID:        0,
							Title:     "test title",
							Text:      "test text",
							UserID:    0,
							Published: true,
						}, nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 0,
						"title": "test title",
						"text": "test text",
						"user_id": 0,
						"published": true,
						"date_creation": "",
						"date_update": ""
					}
				}
				`,
		},
		{
			name: "invalid json body passed",
			adID: "0",
			ad: struct {
				UserID bool `json:"user_id"`
			}{
				UserID: true,
			},
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "json: cannot unmarshal bool into Go struct field ChangeAdStatusRequest.user_id of type int64"}`,
		},
		{
			name:               "invalid user id passed",
			adID:               "invalid_user_id_1",
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_user_id_1\": invalid syntax"}`,
		},
		{
			name: "error from service: ErrNoAccess",
			adID: "0",
			ad: request.ChangeAdStatusRequest{
				UserID:    0,
				Published: true,
			},
			mockBehaviour: func(serv *handlerMock.MockAdService) {
				serv.EXPECT().ChangeAdStatus(gomock.Any(), int64(0), int64(0), true).
					Return(nil, service.ErrNoAccess{Err: service.ErrNoAccessAd})
			},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   `{"error": "you don't have access to edit the adID"}`,
		},
		{
			name: "error from service",
			adID: "0",
			ad: request.ChangeAdStatusRequest{
				Published: true,
				UserID:    0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().ChangeAdStatus(gomock.Any(), int64(0), int64(0), true).
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.PUT("/:ad_id/status", handler.changeAdStatus)

			jsonValue, err := json.Marshal(tc.ad)
			require.Equal(t, err, nil)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/status", tc.adID), bytes.NewBuffer(jsonValue))

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_updateAd(t *testing.T) {
	tests := []struct {
		name               string
		adID               string
		ad                 any
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			adID: "0",
			ad: request.UpdateAdRequest{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					UpdateAd(gomock.Any(), int64(0), int64(0), "test title", "test text").
					Return(
						&models.Ad{
							ID:        0,
							Title:     "test title",
							Text:      "test text",
							UserID:    0,
							Published: true,
						}, nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": {
						"id": 0,
						"title": "test title",
						"text": "test text",
						"user_id": 0,
						"published": true,
						"date_creation": "",
						"date_update": ""
					}
				}
				`,
		},
		{
			name: "invalid json body passed",
			adID: "0",
			ad: struct {
				UserID bool `json:"user_id"`
			}{
				UserID: true,
			},
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "json: cannot unmarshal bool into Go struct field UpdateAdRequest.user_id of type int64"}`,
		},
		{
			name:               "invalid user id passed",
			adID:               "invalid_user_id_1",
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_user_id_1\": invalid syntax"}`,
		},
		{
			name: "error from service: ErrNoAccess",
			adID: "0",
			ad: request.UpdateAdRequest{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(serv *handlerMock.MockAdService) {
				serv.EXPECT().UpdateAd(gomock.Any(), int64(0), int64(0), "test title", "test text").
					Return(nil, service.ErrNoAccess{Err: service.ErrNoAccessAd})
			},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   `{"error": "you don't have access to edit the adID"}`,
		},
		{
			name: "error from service: Validation",
			adID: "0",
			ad: request.UpdateAdRequest{
				Title:  "",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().UpdateAd(gomock.Any(), int64(0), int64(0), "", "test text").
					Return(nil, publication.ValidationErrors{publication.ValidationError{Err: publication.ErrInvalidTitle}})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "wrong title"}`,
		},
		{
			name: "error from service",
			adID: "0",
			ad: request.UpdateAdRequest{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().UpdateAd(gomock.Any(), int64(0), int64(0), "test title", "test text").
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.PUT("/:ad_id", handler.updateAd)

			jsonValue, err := json.Marshal(tc.ad)
			require.Equal(t, err, nil)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s", tc.adID), bytes.NewBuffer(jsonValue))

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_deleteAd(t *testing.T) {
	tests := []struct {
		name               string
		adID               string
		ad                 any
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			adID: "0",
			ad: request.DeleteAdRequest{
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					DeleteAd(gomock.Any(), int64(0), int64(0)).
					Return(
						nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"success":"User #0 deleted"}`,
		},
		{
			name: "invalid json body passed",
			adID: "0",
			ad: struct {
				UserID bool `json:"user_id"`
			}{
				UserID: true,
			},
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "json: cannot unmarshal bool into Go struct field DeleteAdRequest.user_id of type int64"}`,
		},
		{
			name:               "invalid user id passed",
			adID:               "invalid_user_id_1",
			mockBehaviour:      func(service *handlerMock.MockAdService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error": "strconv.Atoi: parsing \"invalid_user_id_1\": invalid syntax"}`,
		},
		{
			name: "error from service: ErrNoAccess",
			adID: "0",
			ad: request.DeleteAdRequest{
				UserID: 0,
			},
			mockBehaviour: func(serv *handlerMock.MockAdService) {
				serv.EXPECT().DeleteAd(gomock.Any(), int64(0), int64(0)).
					Return(service.ErrNoAccess{Err: service.ErrNoAccessAd})
			},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   `{"error": "you don't have access to edit the adID"}`,
		},
		{
			name: "error from service",
			adID: "0",
			ad: request.DeleteAdRequest{
				UserID: 0,
			},
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().DeleteAd(gomock.Any(), int64(0), int64(0)).
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.DELETE("/:ad_id", handler.deleteAd)

			jsonValue, err := json.Marshal(tc.ad)
			require.Equal(t, err, nil)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", tc.adID), bytes.NewBuffer(jsonValue))

			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_searchAds(t *testing.T) {
	tests := []struct {
		name               string
		text               string
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			text: "test",
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					GetAdsByTitle(gomock.Any(), "test").
					Return(
						[]*models.Ad{
							{
								ID:        0,
								Title:     "test title",
								Text:      "test text",
								UserID:    0,
								Published: true,
							},
						},
						nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": [
						{
							"id": 0,
							"title": "test title",
							"text": "test text",
							"user_id": 0,
							"published": true,
							"date_creation": "",
							"date_update": ""
						}
					]
				}
				`,
		},
		{
			name: "error from service",
			text: "test",
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().GetAdsByTitle(gomock.Any(), "test").
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.GET("/search", handler.searchAds)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/search?text=%s", tc.text), nil)
			//r := httptest.NewRequest(http.MethodGet, "/search", nil)
			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserHandler_listAds(t *testing.T) {
	tests := []struct {
		name               string
		text               string
		mockBehaviour      func(service *handlerMock.MockAdService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successfully get user",
			text: "test",
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().
					ListAds(gomock.Any(), "", "", "").
					Return(
						[]*models.Ad{
							{
								ID:        0,
								Title:     "test title",
								Text:      "test text",
								UserID:    0,
								Published: true,
							},
						},
						nil,
					).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: `
				{
					"data": [
						{
							"id": 0,
							"title": "test title",
							"text": "test text",
							"user_id": 0,
							"published": true,
							"date_creation": "",
							"date_update": ""
						}
					]
				}
				`,
		},
		{
			name: "error from service",
			text: "test",
			mockBehaviour: func(service *handlerMock.MockAdService) {
				service.EXPECT().ListAds(gomock.Any(), "", "", "").
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

			service := handlerMock.NewMockAdService(ctrl)
			tc.mockBehaviour(service)

			handler := NewAdHandler(service, nil)

			//Test Server
			rg := gin.New()
			rg.GET("/", handler.listAds)

			//Test request
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			//Perform request
			rg.ServeHTTP(w, r)

			// Assert
			require.Equal(t, tc.expectedStatusCode, w.Code)
			require.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

package mapper

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
	"testing"
	"time"

	"homework10/internal/api/handlers/httpgin/response"
	"homework10/internal/domain/models"

	"github.com/stretchr/testify/require"
)

func TestAdToResponse(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		ad       *models.Ad
		expected response.AdResponse
	}{
		{
			name: "successfully map ad to response",
			ad: &models.Ad{
				ID:           1,
				Title:        "test title",
				Text:         "test text",
				UserID:       2,
				Published:    false,
				DateCreation: now.Format("01-02-2006"),
				DateUpdate:   now.Format("01-02-2006"),
			},
			expected: response.AdResponse{
				ID:           1,
				Title:        "test title",
				Text:         "test text",
				UserID:       2,
				Published:    false,
				DateCreation: now.Format("01-02-2006"),
				DateUpdate:   now.Format("01-02-2006"),
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			actual := AdToResponse(tc.ad)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestAdToSliceResponse(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		ads      []*models.Ad
		expected []response.AdResponse
	}{
		{
			name: "successfully map adSlice to responseSlice",
			ads: []*models.Ad{
				{
					ID:           1,
					Title:        "test title",
					Text:         "test text",
					UserID:       2,
					Published:    false,
					DateCreation: now.Format("01-02-2006"),
					DateUpdate:   now.Format("01-02-2006"),
				},
			},
			expected: []response.AdResponse{
				{
					ID:           1,
					Title:        "test title",
					Text:         "test text",
					UserID:       2,
					Published:    false,
					DateCreation: now.Format("01-02-2006"),
					DateUpdate:   now.Format("01-02-2006"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdToSliceResponse(tt.ads); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AdToSliceResponse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestAdSuccessResponse(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		ad       *models.Ad
		expected *fiber.Map
	}{
		{
			name: "successfully map ad to fiber.Map",
			ad: &models.Ad{
				ID:           1,
				Title:        "test title",
				Text:         "test text",
				UserID:       2,
				Published:    false,
				DateCreation: now.Format("01-02-2006"),
				DateUpdate:   now.Format("01-02-2006"),
			},
			expected: &fiber.Map{
				"data": response.AdResponse{
					ID:           1,
					Title:        "test title",
					Text:         "test text",
					UserID:       2,
					Published:    false,
					DateCreation: now.Format("01-02-2006"),
					DateUpdate:   now.Format("01-02-2006"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdSuccessResponse(tt.ad); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AdSuccessResponse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestAdsSuccessResponse(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		ads      []*models.Ad
		expected *fiber.Map
	}{
		{
			name: "successfully map adSlice to fiber.Map",
			ads: []*models.Ad{
				{
					ID:           1,
					Title:        "test title",
					Text:         "test text",
					UserID:       2,
					Published:    false,
					DateCreation: now.Format("01-02-2006"),
					DateUpdate:   now.Format("01-02-2006"),
				},
			},
			expected: &fiber.Map{
				"data": []response.AdResponse{
					{
						ID:           1,
						Title:        "test title",
						Text:         "test text",
						UserID:       2,
						Published:    false,
						DateCreation: now.Format("01-02-2006"),
						DateUpdate:   now.Format("01-02-2006"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdsSuccessResponse(tt.ads); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AdsSuccessResponse() = %v, want %v", got, tt.expected)
			}
		})
	}
}

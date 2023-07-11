package mapper

import (
	"github.com/stretchr/testify/require"
	contracts "homework10/internal/api/handlers/grpc/contracts/langs/go"
	"homework10/internal/domain/models"
	"reflect"
	"testing"
)

func TestAdToResponse(t *testing.T) {

	tests := []struct {
		name     string
		ad       *models.Ad
		expected *contracts.AdResponse
	}{
		{
			name: "successfully map ad to response",
			ad: &models.Ad{
				ID:        1,
				Title:     "test title",
				Text:      "test text",
				UserID:    2,
				Published: false,
			},
			expected: &contracts.AdResponse{
				Id:        1,
				Title:     "test title",
				Text:      "test text",
				UserId:    2,
				Published: false,
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
	tests := []struct {
		name     string
		ads      []*models.Ad
		expected *contracts.ListAdsResponse
	}{
		{
			name: "successfully map adSlice to response slice",
			ads: []*models.Ad{
				{
					ID:        1,
					Title:     "test title",
					Text:      "test text",
					UserID:    2,
					Published: false,
				},
			},
			expected: &contracts.ListAdsResponse{
				List: []*contracts.AdResponse{
					{
						Id:        1,
						Title:     "test title",
						Text:      "test text",
						UserId:    2,
						Published: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdsToListResponse(tt.ads); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("AdToSliceResponse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

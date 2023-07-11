package localrepo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"homework10/internal/domain/models"
	"testing"
)

func TestAdRepo_GetAd(t *testing.T) {
	adRepo := NewAdRepo()
	adRepo.storage[0] = &models.Ad{
		ID:        0,
		Title:     "test title",
		Text:      "test text",
		UserID:    0,
		Published: false,
	}
	tests := []struct {
		name     string
		adID     int64
		expected *models.Ad
		err      error
		cancel   bool
	}{
		{
			name:     "successfully test AddAd",
			adID:     0,
			expected: adRepo.storage[0],
			cancel:   false,
		},
		{
			name: "successfully test AddAd",
			adID: 10,

			expected: nil,
			err:      fmt.Errorf("the ad does not exist"),
			cancel:   false,
		},
		{
			name:     "error context",
			expected: nil,
			err:      context.Canceled,
			cancel:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			ad, err := adRepo.GetAd(ctx, tc.adID)
			if err != nil {
				assert.Nil(t, ad)
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, *tc.expected, *ad)
			}
		})
	}
}

func TestAdRepo_GetAds(t *testing.T) {
	adRepo := NewAdRepo()

	adRepo.storage[0] = &models.Ad{
		ID:        0,
		Title:     "test title",
		Text:      "test text",
		UserID:    0,
		Published: false,
	}

	tests := []struct {
		name     string
		expected []*models.Ad
		err      error
		cancel   bool
	}{
		{
			name:     "successfully test AddAd",
			expected: []*models.Ad{adRepo.storage[0]},
			cancel:   false,
		},
		{
			name:     "error context",
			expected: nil,
			err:      context.Canceled,
			cancel:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			ads, err := adRepo.GetAds(ctx)
			if err != nil {
				assert.Nil(t, ads)
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, len(tc.expected), len(ads))
				assert.Equal(t, *tc.expected[0], *ads[0])
			}
		})
	}
}

func TestAdRepo_AddAd(t *testing.T) {
	tests := []struct {
		name     string
		adRepo   *AdRepo
		ad       models.Ad
		expected int64
		err      error
		cancel   bool
	}{
		{
			name:   "successfully test AddAd",
			adRepo: NewAdRepo(),
			ad: models.Ad{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			expected: 0,
			cancel:   false,
		},
		{
			name:   "error context",
			adRepo: NewAdRepo(),
			ad: models.Ad{
				Title:  "test title",
				Text:   "test text",
				UserID: 0,
			},
			expected: 0,
			err:      context.Canceled,
			cancel:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			adID, err := tc.adRepo.AddAd(ctx, tc.ad)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, adID, tc.expected)
		})
	}
}

func TestAdRepo_SetStatus(t *testing.T) {
	adRepo := NewAdRepo()

	adRepo.storage[0] = &models.Ad{
		ID:        0,
		Title:     "test title",
		Text:      "test text",
		UserID:    0,
		Published: false,
	}

	tests := []struct {
		name      string
		adID      int64
		expected  *models.Ad
		published bool
		err       error
		cancel    bool
	}{
		{
			name: "successfully test AddAd",
			adID: 0,
			expected: &models.Ad{
				ID:        0,
				Title:     "test title",
				Text:      "test text",
				UserID:    0,
				Published: true,
			},
			published: true,
			cancel:    false,
		},
		{
			name: "successfully test AddAd",
			adID: 0,
			expected: &models.Ad{
				ID:        0,
				Title:     "test title",
				Text:      "test text",
				UserID:    0,
				Published: true,
			},
			published: true,
			err:       context.Canceled,
			cancel:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithCancelCause(context.Background())

			if tc.cancel {
				cancel(context.Canceled)
			}

			ad, err := adRepo.SetStatus(ctx, tc.adID, tc.published)
			if err != nil {
				assert.Nil(t, ad)
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, *tc.expected, *ad)
			}
		})
	}
}

func TestAdRepo_Update(t *testing.T) {
	adRepo := NewAdRepo()

	adRepo.storage[0] = &models.Ad{
		ID:        0,
		Title:     "test title",
		Text:      "test text",
		UserID:    0,
		Published: false,
	}

	tests := []struct {
		name     string
		adID     int64
		expected *models.Ad
		title    string
		text     string
		err      error
		cancel   bool
	}{
		{
			name:  "successfully test AddAd",
			adID:  0,
			title: "new title",
			text:  "new text",
			expected: &models.Ad{
				ID:        0,
				Title:     "new title",
				Text:      "new text",
				UserID:    0,
				Published: false,
			},
			cancel: false,
		},
		{

			name:  "successfully test AddAd",
			adID:  0,
			title: "new title",
			text:  "new text",
			expected: &models.Ad{
				ID:        0,
				Title:     "new title",
				Text:      "new text",
				UserID:    0,
				Published: false,
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

			ad, err := adRepo.Update(ctx, tc.adID, tc.title, tc.text)
			if err != nil {
				assert.Nil(t, ad)
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, *tc.expected, *ad)
			}
		})
	}
}

func TestAdRepo_Delete(t *testing.T) {
	adRepo := NewAdRepo()

	adRepo.storage[0] = &models.Ad{
		ID:        0,
		Title:     "test title",
		Text:      "test text",
		UserID:    0,
		Published: false,
	}

	tests := []struct {
		name   string
		adID   int64
		err    error
		cancel bool
	}{
		{
			name:   "successfully test AddAd",
			adID:   0,
			err:    nil,
			cancel: false,
		},
		{
			name:   "successfully test AddAd",
			adID:   0,
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

			err := adRepo.DeleteAd(ctx, tc.adID)
			if err != nil {
				assert.Equal(t, tc.err, err)
			} else {
				assert.Equal(t, len(adRepo.storage), 0)
			}
		})
	}
}

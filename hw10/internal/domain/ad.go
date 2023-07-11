package domain

import (
	"context"
	"homework10/internal/domain/models"
)

//go:generate mockgen -source=./ad.go -destination=../service/mock/ad.go -package=repoMock AdRepository
type AdRepository interface {
	AddAd(ctx context.Context, ad models.Ad) (int64, error)
	GetAd(ctx context.Context, adID int64) (*models.Ad, error)
	SetStatus(ctx context.Context, adID int64, published bool) (*models.Ad, error)
	Update(ctx context.Context, adID int64, title string, text string) (*models.Ad, error)
	DeleteAd(ctx context.Context, adID int64) error
	GetAds(ctx context.Context) ([]*models.Ad, error)
}

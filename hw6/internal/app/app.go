package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilgizjan1/publication"
	"homework6/internal/ads"
)

type ErrNoAccess struct {
	err error
}

func (e ErrNoAccess) Error() string {
	return fmt.Sprintf("%s", e.err)
}

var ErrNoAccessAd = errors.New("you don't have access to edit the ad")

type App interface {
	CreateAd(ctx context.Context, title string, text string, authorID int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*ads.Ad, error)
}

type Repository interface {
	AddAd(ctx context.Context, ad ads.Ad) (int64, error)
	GetAuthorID(ctx context.Context, adID int64) (int64, error)
	SetStatus(ctx context.Context, adID int64, published bool) (*ads.Ad, error)
	Update(ctx context.Context, adID int64, title string, text string) (*ads.Ad, error)
}

type AdService struct {
	repository Repository
}

func NewApp(repo Repository) *AdService {
	return &AdService{repository: repo}
}

func (a *AdService) CreateAd(ctx context.Context, title string, text string, userID int64) (*ads.Ad, error) {
	ad := ads.Ad{Title: title, Text: text, AuthorID: userID}

	if err := publication.Validate(ad); err != nil {
		return nil, err
	}

	id, err := a.repository.AddAd(ctx, ad)
	if err != nil {
		return nil, fmt.Errorf("adding add: %w", err)
	}

	ad.ID = id

	return &ad, nil
}

func (a *AdService) ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (*ads.Ad, error) {
	authorID, err := a.repository.GetAuthorID(ctx, adID)
	if err != nil {
		return nil, err
	}
	if authorID != userID {
		return nil, ErrNoAccess{err: ErrNoAccessAd}
	}
	ad, err := a.repository.SetStatus(ctx, adID, published)
	if err != nil {
		return nil, fmt.Errorf("setting ad status: %w", err)
	}
	return ad, nil
}

func (a *AdService) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	authorID, err := a.repository.GetAuthorID(ctx, adID)
	if err != nil {
		return nil, err
	}
	if authorID != userID {
		return nil, ErrNoAccess{err: ErrNoAccessAd}
	}
	ad, err := a.repository.Update(ctx, adID, title, text)
	if err != nil {
		return nil, fmt.Errorf("updating add: %w", err)
	}
	if err := publication.Validate(*ad); err != nil {
		return nil, err
	}
	return ad, nil
}

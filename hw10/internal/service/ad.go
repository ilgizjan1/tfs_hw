package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilgizjan1/publication"
	"homework10/internal/domain"
	"homework10/internal/domain/models"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "01-02-2006"
)

type ErrNoAccess struct {
	Err error
}

func (e ErrNoAccess) Error() string {
	return fmt.Sprintf("%s", e.Err)
}

var ErrNoAccessAd = errors.New("you don't have access to edit the adID")

type AdService struct {
	adRepo domain.AdRepository
}

func NewAdService(adRepo domain.AdRepository) *AdService {
	return &AdService{
		adRepo: adRepo,
	}
}

func (s *AdService) GetAdByID(ctx context.Context, adID int64) (*models.Ad, error) {
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (s *AdService) CreateAd(ctx context.Context, title string, text string, userID int64) (*models.Ad, error) {
	ad := models.Ad{Title: title, Text: text, UserID: userID,
		DateCreation: time.Now().UTC().Format(dateFormat), DateUpdate: time.Now().UTC().Format(dateFormat)}

	if err := publication.Validate(ad); err != nil {
		return nil, err
	}

	id, err := s.adRepo.AddAd(ctx, ad)
	if err != nil {
		return nil, fmt.Errorf("adding add: %w", err)
	}

	ad.ID = id

	return &ad, nil
}

func (s *AdService) ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (*models.Ad, error) {
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}
	if ad.UserID != userID {
		return nil, ErrNoAccess{Err: ErrNoAccessAd}
	}
	newAd, err := s.adRepo.SetStatus(ctx, adID, published)
	if err != nil {
		return nil, fmt.Errorf("setting adID status: %w", err)
	}
	newAd.DateUpdate = time.Now().UTC().Format(dateFormat)
	return newAd, nil
}

func (s *AdService) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*models.Ad, error) {
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}
	if ad.UserID != userID {
		return nil, ErrNoAccess{Err: ErrNoAccessAd}
	}
	newAd, err := s.adRepo.Update(ctx, adID, title, text)
	if err != nil {
		return nil, fmt.Errorf("updating add: %w", err)
	}
	newAd.DateUpdate = time.Now().UTC().Format(dateFormat)
	if err := publication.Validate(*newAd); err != nil {
		return nil, err
	}

	return newAd, nil
}

func (s *AdService) DeleteAd(ctx context.Context, adID int64, userID int64) error {
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return err
	}
	if ad.UserID != userID {
		return ErrNoAccess{Err: ErrNoAccessAd}
	}
	return s.adRepo.DeleteAd(ctx, adID)
}

func (s *AdService) GetAdsByTitle(ctx context.Context, text string) ([]*models.Ad, error) {
	adSlice, err := s.adRepo.GetAds(ctx)
	if err != nil {
		return nil, err
	}
	cleanSlice := make([]*models.Ad, 0)
	for _, ad := range adSlice {
		if strings.Contains(ad.Title, text) {
			cleanSlice = append(cleanSlice, ad)
		}
	}
	return cleanSlice, nil
}

func (s *AdService) ListAds(ctx context.Context, publishedRaw string, userIDRaw string, dateCreationRaw string) ([]*models.Ad, error) {
	adSlice, err := s.adRepo.GetAds(ctx)
	if err != nil {
		return nil, err
	}
	funSlice := make([]funCheckAd, 0)
	if publishedRaw != "" {
		published, err := strconv.ParseBool(publishedRaw)
		if err != nil {
			return nil, err
		}
		funSlice = append(funSlice, checkPublished(published))
	} else if userIDRaw == "" && dateCreationRaw == "" {
		funSlice = append(funSlice, checkPublished(true))
	}

	if userIDRaw != "" {
		userID, err := strconv.Atoi(userIDRaw)
		if err != nil {
			return nil, fmt.Errorf("userID validating error")
		}
		funSlice = append(funSlice, checkUserID(int64(userID)))
	}
	if dateCreationRaw != "" {
		date, err := time.Parse(dateFormat, dateCreationRaw)
		if err != nil {
			return nil, fmt.Errorf("dateCreation validating error")
		}
		funSlice = append(funSlice, checkDate(date))
	}
	cleanSlice := make([]*models.Ad, 0)
	for _, ad := range adSlice {
		if !check(*ad, funSlice) {
			continue
		}
		cleanSlice = append(cleanSlice, ad)
	}
	return cleanSlice, nil
}

type funCheckAd func(ad models.Ad) bool

func checkPublished(published bool) funCheckAd {
	return func(ad models.Ad) bool {
		return ad.Published == published
	}
}

func checkUserID(userID int64) funCheckAd {
	return func(ad models.Ad) bool {
		return ad.UserID == userID
	}
}

func checkDate(date time.Time) funCheckAd {
	return func(ad models.Ad) bool {
		return date.Format(dateFormat) == ad.DateCreation
	}
}

func check(ad models.Ad, functions []funCheckAd) bool {
	for _, fun := range functions {
		if !fun(ad) {
			return false
		}
	}
	return true
}

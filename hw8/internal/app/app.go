package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilgizjan1/publication"
	"homework8/internal/ads"
	"homework8/internal/users"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "01-02-2006"
)

type ErrNoAccess struct {
	err error
}

func (e ErrNoAccess) Error() string {
	return fmt.Sprintf("%s", e.err)
}

var ErrNoAccessAd = errors.New("you don't have access to edit the ad")

type App interface {
	GetAdByID(ctx context.Context, adID int64) (*ads.Ad, error)
	CreateAd(ctx context.Context, title string, text string, authorID int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*ads.Ad, error)
	GetAdsByTitle(ctx context.Context, text string) ([]*ads.Ad, error)
	CreateUser(ctx context.Context, nickName string, email string) (*users.User, error)
	UpdateUser(ctx context.Context, userID int64, nickName string, email string) (*users.User, error)
	ListAds(ctx context.Context, published string, userIDRaw string, dateCreationRaw string) ([]*ads.Ad, error)
}

type AdRepository interface {
	AddAd(ctx context.Context, ad ads.Ad) (int64, error)
	GetAd(ctx context.Context, adID int64) (*ads.Ad, error)
	SetStatus(ctx context.Context, adID int64, published bool) (*ads.Ad, error)
	Update(ctx context.Context, adID int64, title string, text string) (*ads.Ad, error)
	GetAds(ctx context.Context) ([]*ads.Ad, error)
}

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*users.User, error)
	AddUser(ctx context.Context, user users.User) (int64, error)
	Update(ctx context.Context, userID int64, nickName string, email string) (*users.User, error)
}

type Service struct {
	adRepo   AdRepository
	userRepo UserRepository
}

func NewApp(adRepo AdRepository, userRepo UserRepository) *Service {
	return &Service{adRepo: adRepo, userRepo: userRepo}
}

func (s *Service) CreateAd(ctx context.Context, title string, text string, userID int64) (*ads.Ad, error) {
	if _, err := s.userRepo.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	ad := ads.Ad{Title: title, Text: text, AuthorID: userID,
		DateCreation: time.Now().UTC(), DateUpdate: time.Now().UTC()}

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

func (s *Service) ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (*ads.Ad, error) {
	if _, err := s.userRepo.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != userID {
		return nil, ErrNoAccess{err: ErrNoAccessAd}
	}
	newAd, err := s.adRepo.SetStatus(ctx, adID, published)
	newAd.DateUpdate = time.Now().UTC()
	if err != nil {
		return nil, fmt.Errorf("setting ad status: %w", err)
	}
	return newAd, nil
}

func (s *Service) UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*ads.Ad, error) {
	if _, err := s.userRepo.GetUser(ctx, userID); err != nil {
		return nil, err
	}
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != userID {
		return nil, ErrNoAccess{err: ErrNoAccessAd}
	}
	newAd, err := s.adRepo.Update(ctx, adID, title, text)
	if err != nil {
		return nil, fmt.Errorf("updating add: %w", err)
	}
	if err := publication.Validate(*ad); err != nil {
		return nil, err
	}
	return newAd, nil
}

func (s *Service) GetAdByID(ctx context.Context, adID int64) (*ads.Ad, error) {
	ad, err := s.adRepo.GetAd(ctx, adID)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (s *Service) GetAdsByTitle(ctx context.Context, text string) ([]*ads.Ad, error) {
	adSlice, err := s.adRepo.GetAds(ctx)
	if err != nil {
		return nil, err
	}
	cleanSlice := make([]*ads.Ad, 0)
	for _, ad := range adSlice {
		if strings.Contains(ad.Title, text) {
			cleanSlice = append(cleanSlice, ad)
		}
	}
	return cleanSlice, nil
}

type funCheckAd func(ad ads.Ad) bool

func checkPublished(published bool) funCheckAd {
	return func(ad ads.Ad) bool {
		return ad.Published == published
	}
}

func checkUserID(userID int64) funCheckAd {
	return func(ad ads.Ad) bool {
		return ad.AuthorID == userID
	}
}

func checkDate(date time.Time) funCheckAd {
	return func(ad ads.Ad) bool {
		return date.Format(dateFormat) == ad.DateCreation.Format(dateFormat)
	}
}

func check(ad ads.Ad, functions []funCheckAd) bool {
	for _, fun := range functions {
		if !fun(ad) {
			return false
		}
	}
	return true
}

func (s *Service) ListAds(ctx context.Context, publishedRaw string, userIDRaw string, dateCreationRaw string) ([]*ads.Ad, error) {
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
	cleanSlice := make([]*ads.Ad, 0)
	for _, ad := range adSlice {
		if !check(*ad, funSlice) {
			continue
		}
		cleanSlice = append(cleanSlice, ad)
	}
	return cleanSlice, nil
}

func (s *Service) CreateUser(ctx context.Context, nickName string, email string) (*users.User, error) {
	user := users.User{NickName: nickName, Email: email}
	userID, err := s.userRepo.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID
	return &user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID int64, nickName string, email string) (*users.User, error) {
	user, err := s.userRepo.Update(ctx, userID, nickName, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

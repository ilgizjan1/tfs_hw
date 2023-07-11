package grpc

import (
	"context"
	contracts "homework10/internal/api/handlers/grpc/contracts/langs/go"
	"homework10/internal/api/handlers/grpc/mapper"
	"homework10/internal/domain/models"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AdService interface {
	GetAdByID(ctx context.Context, adID int64) (*models.Ad, error)
	CreateAd(ctx context.Context, title string, text string, authorID int64) (*models.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, userID int64, published bool) (*models.Ad, error)
	UpdateAd(ctx context.Context, adID int64, userID int64, title string, text string) (*models.Ad, error)
	DeleteAd(ctx context.Context, adID int64, userID int64) error
	GetAdsByTitle(ctx context.Context, text string) ([]*models.Ad, error)
	ListAds(ctx context.Context, published string, userIDRaw string, dateCreationRaw string) ([]*models.Ad, error)
}

type AdHandler struct {
	adService AdService
}

func NewAdHandler(adServ AdService) *AdHandler {
	return &AdHandler{
		adService: adServ,
	}
}

func (g *AdHandler) GetAd(ctx context.Context, request *contracts.GetAdRequest) (*contracts.AdResponse, error) {
	ad, err := g.adService.GetAdByID(ctx, request.AdId)
	if err != nil {
		return nil, err
	}
	return mapper.AdToResponse(ad), nil
}

func (g *AdHandler) CreateAd(ctx context.Context, request *contracts.CreateAdRequest) (*contracts.AdResponse, error) {
	ad, err := g.adService.CreateAd(ctx, request.Title, request.Text, request.UserId)
	if err != nil {
		return nil, err
	}
	return mapper.AdToResponse(ad), nil
}

func (g *AdHandler) ChangeAdStatus(ctx context.Context, request *contracts.ChangeAdStatusRequest) (*contracts.AdResponse, error) {
	ad, err := g.adService.ChangeAdStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, err
	}
	return mapper.AdToResponse(ad), nil
}

func (g *AdHandler) UpdateAd(ctx context.Context, request *contracts.UpdateAdRequest) (*contracts.AdResponse, error) {
	ad, err := g.adService.UpdateAd(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, err
	}
	return mapper.AdToResponse(ad), nil
}

func (g *AdHandler) DeleteAd(ctx context.Context, request *contracts.DeleteAdRequest) (*emptypb.Empty, error) {
	err := g.adService.DeleteAd(ctx, request.AdId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *AdHandler) SearchAds(ctx context.Context, request *contracts.SearchAdsRequest) (*contracts.ListAdsResponse, error) {
	ads, err := g.adService.GetAdsByTitle(ctx, request.Text)
	if err != nil {
		return nil, err
	}
	return mapper.AdsToListResponse(ads), nil
}

func (g *AdHandler) ListAds(ctx context.Context, request *contracts.ListAdsRequest) (*contracts.ListAdsResponse, error) {
	ads, err := g.adService.ListAds(ctx, request.Published, request.UserId, request.Date)
	if err != nil {
		return nil, err
	}
	return mapper.AdsToListResponse(ads), nil
}

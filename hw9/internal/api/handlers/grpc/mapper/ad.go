package mapper

import (
	contracts "homework9/internal/api/handlers/grpc/contracts/langs/go"
	"homework9/internal/domain/models"
)

func AdToResponse(ad *models.Ad) *contracts.AdResponse {
	return &contracts.AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		UserId:    ad.UserID,
		Published: ad.Published,
	}
}

func AdsToListResponse(ads []*models.Ad) *contracts.ListAdsResponse {
	listAds := make([]*contracts.AdResponse, 0)
	for _, ad := range ads {
		listAds = append(listAds, AdToResponse(ad))
	}
	return &contracts.ListAdsResponse{List: listAds}
}

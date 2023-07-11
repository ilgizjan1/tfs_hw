package mapper

import (
	"github.com/gofiber/fiber/v2"
	"homework10/internal/api/handlers/httpgin/response"
	"homework10/internal/domain/models"
)

func AdToResponse(ad *models.Ad) response.AdResponse {
	return response.AdResponse{
		ID:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		UserID:       ad.UserID,
		Published:    ad.Published,
		DateCreation: ad.DateCreation,
		DateUpdate:   ad.DateUpdate,
	}
}

func AdToSliceResponse(ads []*models.Ad) []response.AdResponse {
	adsRes := make([]response.AdResponse, 0)
	for _, ad := range ads {
		adsRes = append(adsRes, AdToResponse(ad))
	}
	return adsRes
}

func AdSuccessResponse(ad *models.Ad) *fiber.Map {
	return &fiber.Map{
		"data": AdToResponse(ad),
	}
}

func AdsSuccessResponse(ads []*models.Ad) *fiber.Map {
	return &fiber.Map{
		"data": AdToSliceResponse(ads),
	}
}

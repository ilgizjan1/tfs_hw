package mapper

import (
	"github.com/gofiber/fiber/v2"
	"homework9/internal/api/handlers/httpgin/response"
	"homework9/internal/domain/models"
)

func AdToResponse(ad *models.Ad) response.AdResponse {
	return response.AdResponse{
		ID:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorID:  ad.UserID,
		Published: ad.Published,
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
		"data":  AdToResponse(ad),
		"error": nil,
	}
}

func AdsSuccessResponse(ads []*models.Ad) *fiber.Map {
	return &fiber.Map{
		"data":  AdToSliceResponse(ads),
		"error": nil,
	}
}

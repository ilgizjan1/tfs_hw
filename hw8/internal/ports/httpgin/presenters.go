package httpgin

import (
	"github.com/gofiber/fiber/v2"
	"homework8/internal/ads"
	"homework8/internal/users"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type createUserRequest struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

type updateUserRequest struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

type adResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func adToResponse(ad *ads.Ad) adResponse {
	return adResponse{
		ID:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorID:  ad.AuthorID,
		Published: ad.Published,
	}
}

func adToSliceResponse(ads []*ads.Ad) []adResponse {
	adsRes := make([]adResponse, 0)
	for _, ad := range ads {
		adsRes = append(adsRes, adToResponse(ad))
	}
	return adsRes
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func AdSuccessResponse(ad *ads.Ad) *fiber.Map {
	return &fiber.Map{
		//"data": adResponse{
		//	ID:        ad.ID,
		//	Title:     ad.Title,
		//	Text:      ad.Text,
		//	AuthorID:  ad.AuthorID,
		//	Published: ad.Published,
		//},
		"data":  adToResponse(ad),
		"error": nil,
	}
}

func UserSuccessResponse(user *users.User) *fiber.Map {
	return &fiber.Map{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.NickName,
			Email:    user.Email,
		},
		"error": nil,
	}
}

func AdsSuccesResponse(ads []*ads.Ad) *fiber.Map {
	return &fiber.Map{
		"data":  adToSliceResponse(ads),
		"error": nil,
	}
}

type AdFiltersRequest struct {
	Published    bool   `json:"published"`
	AuthorID     string `json:"authorID"`
	DateCreation string `json:"date_creation"`
}

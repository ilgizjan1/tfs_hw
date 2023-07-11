package httpgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ilgizjan1/publication"
	"homework9/internal/api/handlers/httpgin/mapper"
	"homework9/internal/api/handlers/httpgin/middlewares"
	"homework9/internal/api/handlers/httpgin/request"
	"homework9/internal/domain/models"
	"homework9/internal/service"
	"net/http"
	"strconv"
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
	service      AdService
	userIdentity middlewares.UserIdentity
}

func NewAdHandler(service AdService, userIdentity middlewares.UserIdentity) *AdHandler {
	return &AdHandler{
		service:      service,
		userIdentity: userIdentity,
	}
}

func (h *AdHandler) AddRoutes(rg *gin.RouterGroup) {
	rg.GET("/:ad_id", h.getAd)                                                          // Метод для получения объявления (ad) по ID (ad_id)
	rg.POST("/", h.userIdentity.UserIdentityMiddleware(), h.createAd)                   // Метод для создания объявления (ad)
	rg.PUT("/:ad_id/status", h.userIdentity.UserIdentityMiddleware(), h.changeAdStatus) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	rg.PUT("/:ad_id", h.userIdentity.UserIdentityMiddleware(), h.updateAd)              // Метод для обновления текста(Text) или заголовка(Title) объявления
	rg.DELETE("/:ad_id", h.userIdentity.UserIdentityMiddleware(), h.deleteAd)           // Метод для удаления объявления (ad) по ID (ad_id)
	rg.GET("/search", h.searchAds)                                                      // Метод для поиска объявлений по названию (title = "...")
	rg.GET("/", h.listAds)                                                              // Метод для получение списка объявлений с фильтрами
}

func (h *AdHandler) BasePrefix() string {
	return "/ads"
}

// Метод для получения объявления (ad)
func (h *AdHandler) getAd(ctx *gin.Context) {
	adIDRaw := ctx.Param("ad_id")
	adID, err := strconv.Atoi(adIDRaw)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ad, err := h.service.GetAdByID(ctx, int64(adID))
	if err != nil {
		switch err.(type) {
		case service.ErrNoAccess:
			ctx.Status(http.StatusForbidden)
			return
		default:
			ctx.Status(http.StatusBadRequest)
			return
		}
	}
	ctx.IndentedJSON(http.StatusOK, mapper.AdSuccessResponse(ad))
}

// Метод для создания объявления (ad)
func (h *AdHandler) createAd(ctx *gin.Context) {
	var reqBody request.CreateAdRequest
	if err := ctx.ShouldBindBodyWith(&reqBody, binding.JSON); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ad, err := h.service.CreateAd(ctx, reqBody.Title, reqBody.Text, reqBody.UserID)
	if err != nil {
		switch err.(type) {
		case publication.ValidationErrors:
			ctx.Status(http.StatusBadRequest)
			return
		default:
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
	ctx.IndentedJSON(http.StatusOK, mapper.AdSuccessResponse(ad))
}

// Метод для изменения статуса объявления
func (h *AdHandler) changeAdStatus(ctx *gin.Context) {
	var reqBody request.ChangeAdStatusRequest
	if err := ctx.ShouldBindBodyWith(&reqBody, binding.JSON); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	adIDString := ctx.Param("ad_id")
	adID, err := strconv.Atoi(adIDString)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ad, err := h.service.ChangeAdStatus(ctx, int64(adID), reqBody.UserID, reqBody.Published)
	if err != nil {
		switch err.(type) {
		case service.ErrNoAccess:
			ctx.Status(http.StatusForbidden)
			return
		default:
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
	ctx.IndentedJSON(http.StatusOK, mapper.AdSuccessResponse(ad))
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func (h *AdHandler) updateAd(ctx *gin.Context) {
	var reqBody request.UpdateAdRequest
	if err := ctx.ShouldBindBodyWith(&reqBody, binding.JSON); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	adIDRaw := ctx.Param("ad_id")
	adID, err := strconv.Atoi(adIDRaw)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ad, err := h.service.UpdateAd(ctx, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
	if err != nil {
		switch err.(type) {
		case service.ErrNoAccess:
			ctx.Status(http.StatusForbidden)
			return
		case publication.ValidationErrors:
			ctx.Status(http.StatusBadRequest)
			return
		default:
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
	ctx.IndentedJSON(http.StatusOK, mapper.AdSuccessResponse(ad))
}

// Метод для удаления объявления (ad)
func (h *AdHandler) deleteAd(ctx *gin.Context) {
	var reqBody request.DeleteAdRequest
	if err := ctx.ShouldBindBodyWith(&reqBody, binding.JSON); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	adIDRaw := ctx.Param("ad_id")
	adID, err := strconv.Atoi(adIDRaw)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = h.service.DeleteAd(ctx, int64(adID), reqBody.UserID)
	if err != nil {
		switch err.(type) {
		case service.ErrNoAccess:
			ctx.Status(http.StatusForbidden)
			return
		default:
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}
	//ctx.IndentedJSON(http.StatusOK, AdSuccessResponse(ad))
	ctx.IndentedJSON(http.StatusOK, gin.H{"success": "User #" + adIDRaw + " deleted"})
}

// Метод для поиска объявлений (ads) по названию (title)
func (h *AdHandler) searchAds(ctx *gin.Context) {
	text := ctx.Query("text")
	adList, err := h.service.GetAdsByTitle(ctx, text)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.IndentedJSON(http.StatusOK, mapper.AdsSuccessResponse(adList))
}

// Метод для получения объявлений (ads) с возможностью фильтрации
func (h *AdHandler) listAds(ctx *gin.Context) {
	published := ctx.Query("published")
	userID := ctx.Query("user_id")
	dateCreation := ctx.Query("date")
	adSlice, err := h.service.ListAds(ctx, published, userID, dateCreation)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.IndentedJSON(http.StatusOK, mapper.AdsSuccessResponse(adSlice))
}

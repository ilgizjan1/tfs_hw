package httpgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/ilgizjan1/publication"
	"homework8/internal/app"
)

type AdHandler struct {
	app app.App
}

func New(app app.App) *AdHandler {
	return &AdHandler{app: app}
}

// Метод для создания объявления (ad)
func createAd(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody createAdRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ad, err := a.app.CreateAd(ctx, reqBody.Title, reqBody.Text, reqBody.UserID)
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
		ctx.IndentedJSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления
func changeAdStatus(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		adIDString := ctx.Param("ad_id")
		adID, err := strconv.Atoi(adIDString)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ad, err := a.app.ChangeAdStatus(ctx, int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			switch err.(type) {
			case app.ErrNoAccess:
				ctx.Status(http.StatusForbidden)
				return
			default:
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}
		ctx.IndentedJSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody updateAdRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		adIDString := ctx.Param("ad_id")
		adID, err := strconv.Atoi(adIDString)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ad, err := a.app.UpdateAd(ctx, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			switch err.(type) {
			case app.ErrNoAccess:
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
		ctx.IndentedJSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func getAdByID(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		adIDRaw := ctx.Param("ad_id")
		adID, err := strconv.Atoi(adIDRaw)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ad, err := a.app.GetAdByID(ctx, int64(adID))
		if err != nil {
			switch err.(type) {
			case app.ErrNoAccess:
				ctx.Status(http.StatusForbidden)
				return
			default:
				ctx.Status(http.StatusBadRequest)
				return
			}
		}
		ctx.IndentedJSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

func searchAdsByTitle(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		text := ctx.Query("text")
		ads, err := a.app.GetAdsByTitle(ctx, text)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.IndentedJSON(http.StatusOK, AdsSuccesResponse(ads))
	}
}

func listAds(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		published := ctx.Query("published")
		userID := ctx.Query("user_id")
		dateCreation := ctx.Query("date")
		adSlice, err := a.app.ListAds(ctx, published, userID, dateCreation)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.IndentedJSON(http.StatusOK, AdsSuccesResponse(adSlice))
	}
}

func createUser(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody createUserRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		user, err := a.app.CreateUser(ctx, reqBody.NickName, reqBody.Email)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.IndentedJSON(http.StatusOK, UserSuccessResponse(user))
	}
}

func updateUser(a *AdHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody updateUserRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		user, err := a.app.UpdateUser(ctx, reqBody.ID, reqBody.NickName, reqBody.Email)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.IndentedJSON(http.StatusOK, UserSuccessResponse(user))
	}
}

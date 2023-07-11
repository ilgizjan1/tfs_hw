package httpfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ilgizjan1/publication"
	"homework6/internal/app"
)

type AdHandler struct {
	app app.App
}

func New(app app.App) *AdHandler {
	return &AdHandler{app: app}
}

// Метод для создания объявления (ad)
func createAd(a *AdHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.app.CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			switch err.(type) {
			case publication.ValidationErrors:
				c.Status(http.StatusBadRequest)
				return c.JSON(AdErrorResponse(err))
			default:
				c.Status(http.StatusInternalServerError)
				return c.JSON(AdErrorResponse(err))
			}
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления
func changeAdStatus(a *AdHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.app.ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)
		if err != nil {
			switch err.(type) {
			case app.ErrNoAccess:
				c.Status(http.StatusForbidden)
				return c.JSON(AdErrorResponse(err))
			default:
				c.Status(http.StatusInternalServerError)
				return c.JSON(AdErrorResponse(err))
			}
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a *AdHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.app.UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			switch err.(type) {
			case app.ErrNoAccess:
				c.Status(http.StatusForbidden)
				return c.JSON(AdErrorResponse(err))
			case publication.ValidationErrors:
				c.Status(http.StatusBadRequest)
				return c.JSON(AdErrorResponse(err))
			default:
				c.Status(http.StatusInternalServerError)
				return c.JSON(AdErrorResponse(err))
			}
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

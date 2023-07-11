package httpfiber

import (
	"github.com/gofiber/fiber/v2"

	"homework6/internal/app"
)

func AppRouter(r fiber.Router, a app.App) {
	addHandler := New(a)
	r.Post("/ads", createAd(addHandler))                    // Метод для создания объявления (ad)
	r.Put("/ads/:ad_id/status", changeAdStatus(addHandler)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.Put("/ads/:ad_id", updateAd(addHandler))              // Метод для обновления текста(Text) или заголовка(Title) объявления
}

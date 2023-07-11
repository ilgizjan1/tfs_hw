package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	addHandler := New(a)

	r.POST("/ads", createAd(addHandler))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(addHandler)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(addHandler))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAdByID(addHandler))             // Метод для получения объявления (ad) по ID (ad_id)
	r.GET("/ads/search", searchAdsByTitle(addHandler))      // Метод для поиска объявлений по названию (title = "...")
	r.GET("/ads", listAds(addHandler))                      // Метод для получение списка объявлений с фильтрами: только опубликованные, по автору, по дате создания
	r.POST("/users", createUser(addHandler))                // Метод для создания пользователя (user)
	r.PUT("/users/:user_id", updateUser(addHandler))        // Метод для обновления никнейма (Nickname) или почты (Email) пользователя
}

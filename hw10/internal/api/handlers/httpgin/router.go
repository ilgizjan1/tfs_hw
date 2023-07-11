package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/internal/api/handlers/httpgin/middlewares"
)

type ApiVersion string

const (
	ApiV1 ApiVersion = "api/v1"
)

type Router interface {
	AddRoutes(rg *gin.RouterGroup)
	BasePrefix() string
}

func MakeRoutes(apiVersion ApiVersion, routers ...Router) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(
		middlewares.LoggingMiddleware(),
		middlewares.RecoverMiddleware(),
	)

	apiVersionGroup := r.Group(string(apiVersion))
	for _, router := range routers {
		router.AddRoutes(apiVersionGroup.Group(router.BasePrefix()))
	}

	return r
}

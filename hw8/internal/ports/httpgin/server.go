package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
	"homework8/internal/middlewares"
	"net/http"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}
	s.app.Use(middlewares.LoggingMiddleware())
	s.app.Use(gin.Recovery())
	api := s.app.Group("/api/v1")

	AppRouter(api, a)
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}

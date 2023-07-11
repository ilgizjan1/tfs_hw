package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func RecoverMiddleware() gin.HandlerFunc {
	logger := log.New()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Printf("[Recovery] %s panic recovered: %s\n", time.Now(), err)
			}
		}()
		c.Next()
	}
}

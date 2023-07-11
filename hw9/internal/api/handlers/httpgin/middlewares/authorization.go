package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"homework9/internal/domain/models"
	"net/http"
)

type userIDRequest struct {
	UserID int64 `json:"user_id"`
}

type HTTPUserService interface {
	GetUser(ctx context.Context, userID int64) (*models.User, error)
}

type UserIdentity interface {
	UserIdentityMiddleware() gin.HandlerFunc
}

type UserIdentityMiddleware struct {
	service HTTPUserService
}

func NewUserIdentityMiddleware(s HTTPUserService) *UserIdentityMiddleware {
	return &UserIdentityMiddleware{
		service: s,
	}
}

func (a *UserIdentityMiddleware) UserIdentityMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body userIDRequest
		if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Abort()
			return
		}
		_, err := a.service.GetUser(ctx, body.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "the user is not registered"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

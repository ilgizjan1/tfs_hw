package httpgin

import (
	"context"
	"net/http"
	"strconv"

	"homework10/internal/api/handlers/httpgin/mapper"
	"homework10/internal/api/handlers/httpgin/request"
	"homework10/internal/domain/models"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=./user.go -destination=../mock/user.go -package=handlermock UserService
type UserService interface {
	CreateUser(ctx context.Context, nickName string, email string) (*models.User, error)
	UpdateUser(ctx context.Context, userID int64, nickName string, email string) (*models.User, error)
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	DeleteUser(ctx context.Context, userID int64) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) AddRoutes(rg *gin.RouterGroup) {
	rg.GET("/:user_id", h.getUser)       // Метод для получения пользователя (user) по ID (user_id)
	rg.POST("", h.createUser)            // Метод для создания пользователя (user)
	rg.PUT("/:user_id", h.updateUser)    // Метод для обновления никнейма (Nickname) или почты (Email) пользователя
	rg.DELETE("/:user_id", h.deleteUser) // Метод для удаления пользователя (user) по его ID (user_id)и
}

func (h *UserHandler) BasePrefix() string {
	return "/users"
}

func (h *UserHandler) getUser(ctx *gin.Context) {
	userIDRaw := ctx.Param("user_id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}
	user, err := h.service.GetUser(ctx, int64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrResponse(err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, mapper.UserSuccessResponse(user))
}

func (h *UserHandler) createUser(ctx *gin.Context) {
	var reqBody request.CreateUserRequest
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}
	user, err := h.service.CreateUser(ctx, reqBody.NickName, reqBody.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrResponse(err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, mapper.UserSuccessResponse(user))
}

func (h *UserHandler) updateUser(ctx *gin.Context) {
	var reqBody request.UpdateUserRequest
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}
	userIDRaw := ctx.Param("user_id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}
	user, err := h.service.UpdateUser(ctx, int64(userID), reqBody.NickName, reqBody.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrResponse(err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, mapper.UserSuccessResponse(user))
}

func (h *UserHandler) deleteUser(ctx *gin.Context) {
	userIDRaw := ctx.Param("user_id")
	userID, err := strconv.Atoi(userIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrResponse(err))
		return
	}
	err = h.service.DeleteUser(ctx, int64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrResponse(err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"success": "User #" + userIDRaw + " deleted"})
}

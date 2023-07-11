package interceptors

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework9/internal/domain/models"
)

type userIDRequest struct {
	UserID int64 `json:"user_id"`
}

type UserService interface {
	GetUser(ctx context.Context, userID int64) (*models.User, error)
}

type GRPCUserIdentityMiddleware struct {
	service UserService
}

func NewGRPCUserIdentityMiddleware(s UserService) *GRPCUserIdentityMiddleware {
	return &GRPCUserIdentityMiddleware{
		service: s,
	}
}

func (h *GRPCUserIdentityMiddleware) GRPCUserMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	switch info.FullMethod {
	case "/service.AdService/CreateAd",
		"/service.AdService/ChangeAdStatus",
		"/service.AdService/UpdateAd",
		"/service.AdService/DeleteAd":
		a, err := json.Marshal(req)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
		var reqBody userIDRequest
		err = json.Unmarshal(a, &reqBody)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
		_, err = h.service.GetUser(ctx, reqBody.UserID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "the user is not registered")
		}
	}
	return handler(ctx, req)
}

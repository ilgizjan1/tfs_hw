package interceptors

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RecoverInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	logger := log.New()

	defer func() {
		if err := recover(); err != nil {
			logger.Printf("[Recovery] %s panic recovered from method %s: %s\n", time.Now(), info.FullMethod, err)
		}
	}()

	resp, err := handler(ctx, req)
	return resp, err
}

package interceptors

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func LoggingInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	st, _ := status.FromError(err)

	log.WithFields(log.Fields{
		"METHOD":  info.FullMethod,
		"STATUS":  st.Code(),
		"LATENCY": time.Since(start),
		"Error":   err,
	}).Info("GRPC REQUEST")
	return h, err
}

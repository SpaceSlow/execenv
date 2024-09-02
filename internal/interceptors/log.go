package interceptors

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/SpaceSlow/execenv/internal/logger"
)

func size(response interface{}) int {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(response)
	if err != nil {
		return 0
	}
	return binary.Size(buff.Bytes())
}

func LogUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	response, err := handler(ctx, req)

	duration := time.Since(start)

	logger.Log.Info(
		"request/response",
		zap.String("grpc method", info.FullMethod),
		zap.Duration("duration", duration),
		zap.Any("status", status.Code(err)),
		zap.Int("size", size(response)),
	)

	return response, err
}

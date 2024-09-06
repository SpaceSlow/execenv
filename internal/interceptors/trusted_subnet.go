package interceptors

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/SpaceSlow/execenv/internal/config"
)

func WithCheckingTrustedSubnetUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	cfg, err := config.GetServerConfig()
	if err != nil {
		return nil, err
	}

	if cfg.TrustedSubnet != config.NewCIDR("") {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "")
		}
		if realIP := md.Get("X-Real-IP"); len(realIP) != 1 || !cfg.TrustedSubnet.Contains(net.ParseIP(realIP[0])) {
			return nil, status.Error(codes.PermissionDenied, "")
		}
	}

	return handler(ctx, req)
}

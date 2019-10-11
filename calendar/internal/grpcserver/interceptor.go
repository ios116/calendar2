package grpcserver

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func newInterceptor(logger *zap.Logger, token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, errMissingMetadata
		}

		if !valid(md["authorization"], token) {
			return nil, errInvalidToken
		}

		m, err := handler(ctx, req)
		if err != nil {
			logger.Info("RPC failed with error", zap.String("err", err.Error()))
		}

		return m, err
	}
}

// valid validates the authorization.
func valid(authorization []string, token string) bool {
	if len(authorization) < 1 {
		return false
	}
	return token == strings.TrimPrefix(authorization[0], "Bearer ")
}

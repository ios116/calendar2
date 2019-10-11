package web

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ios116/calendar/internal/config"
	gw "github.com/ios116/calendar/internal/grpcserver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": t.Token,
	}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	return false
}

type HttpServer struct {
	HttConfig  *config.HttpConf
	GrpcConfig *config.GrpcConf
	Logger     *zap.Logger
}

func NewHttpServer(httConfig *config.HttpConf, grpcConfig *config.GrpcConf, logger *zap.Logger) *HttpServer {
	return &HttpServer{HttConfig: httConfig, GrpcConfig: grpcConfig, Logger: logger}
}

func (s *HttpServer) newRouter(ctx context.Context) *http.ServeMux {
	// GRPC gate way
	addressRpc := fmt.Sprintf("%s:%d", s.GrpcConfig.GrpcHost, s.GrpcConfig.GrpcPort)
	option := grpc.WithPerRPCCredentials(&tokenAuth{"Bearer secret"})
	conn, err := grpc.Dial(addressRpc, option, grpc.WithInsecure())
	if err != nil {
		s.Logger.Fatal(err.Error())
	}
	rpcGWMux := runtime.NewServeMux()
	err = gw.RegisterCalendarHandler(ctx, rpcGWMux, conn)
	if err != nil {
		s.Logger.Fatal(err.Error())
	}
	// Http
	mux := http.NewServeMux()
	mux.Handle("/v1/", rpcGWMux)
	return mux
}

func (s *HttpServer) Run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dsn := fmt.Sprintf("%s:%d", s.HttConfig.Host, s.HttConfig.Port)
	httpServer := http.Server{
		Addr:    dsn,
		Handler: s.newRouter(ctx),
	}
	if err := httpServer.ListenAndServe(); err != nil {
		s.Logger.Fatal(err.Error())
	}
}

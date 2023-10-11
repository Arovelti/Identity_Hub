package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	proto "github.com/Arovelti/identityhub/api/proto"
	"github.com/Arovelti/identityhub/logger"
	"github.com/Arovelti/identityhub/profile_service/auth"
	"github.com/Arovelti/identityhub/profile_service/rpc"
	"github.com/Arovelti/identityhub/repository"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sagikazarmark/slog-shim"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l := logger.InitLogger()

	// Envs
	host := os.Getenv("LOCAL_HOST")
	portGRPC := os.Getenv("LOCAL_GRPC_PORT")
	portHTTP := os.Getenv("LOCAL_HTTP_PORT")
	localhostGRPC := fmt.Sprintf("%s:%s", host, portGRPC)
	// localhostHTTP := fmt.Sprintf("%s:%s", host, portHTTP)

	network := os.Getenv("NETWORK_TCP")

	// In-Memory Repository
	repo := repository.New()
	repo.GenerateTestProfiles()

	// An instance of the RPC struct with the repository
	rpc := rpc.NewRPC(repo)

	// gRPC server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.BasicAuthInterceptor(repo)),
	)
	proto.RegisterProfilesServer(server, rpc)

	listener, err := net.Listen(network, ":"+portGRPC)
	if err != nil {
		l.Error("Failed to listen gRPC", slog.String("gRPC_listener", err.Error()))
	}
	l.Info("gRPC server listening", slog.String("port", localhostGRPC))

	go func() {
		if err := server.Serve(listener); err != nil {
			l.Error("Failed to serve gRPC", slog.String("gRPC_serve", err.Error()))
		}
	}()

	// Rest HTTP
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = proto.RegisterProfilesHandlerFromEndpoint(context.Background(), mux, localhostGRPC, opts)
	if err != nil {
		l.Error("Failed to register HTTP handler", slog.String("http_handler", err.Error()))
	}

	httpServer := &http.Server{
		Addr:    ":" + portHTTP,
		Handler: mux,
	}
	l.Info("HTTP server listening", slog.String("port", portHTTP))

	if err := httpServer.ListenAndServe(); err != nil {
		l.Error("HTTP server failed to start", slog.String("http_server", err.Error()))
	}
}

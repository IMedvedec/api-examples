package grpc

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/imedvedec/api-examples/api"
	"github.com/imedvedec/api-examples/api/grpc/build"
	"github.com/imedvedec/api-examples/api/grpc/greeting"
	"github.com/imedvedec/api-examples/cert"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type grpcServer struct {
	logger *zerolog.Logger
	server *http.Server
}

func New(addr string) api.Server {
	ctx := context.Background()

	gs := grpc.NewServer()
	build.RegisterGreetingServer(gs, greeting.New())

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := build.RegisterGreetingHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		log.Fatalf("api/grpc/New: greeting handler registration has failed, %v", err)
	}

	tlsCertificate, err := tls.X509KeyPair(cert.ServerCertificate, cert.ServerPrivateKey)
	if err != nil {
		log.Fatalf("api/grpc/New: tls cetificate setup has failed, %v", err)
	}

	hs := http.Server{
		Addr:    addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCertificate},
		},
	}

	consoleLogger := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleLogger).With().Timestamp().Logger()

	grpcServer := grpcServer{
		logger: &logger,
		server: &hs,
	}

	return &grpcServer
}

func (gs *grpcServer) ListenAndServe() error {
	return gs.server.ListenAndServeTLS("", "")
}

func (gs *grpcServer) Shutdown(ctx context.Context) error {
	return gs.server.Shutdown(ctx)
}
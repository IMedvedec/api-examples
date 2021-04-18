package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/imedvedec/api-examples/api"
	"github.com/imedvedec/api-examples/api/grpc/build"
	"github.com/imedvedec/api-examples/api/grpc/greeting"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type grpcServer struct {
	logger     *zerolog.Logger
	restServer *http.Server
	grpcServer *grpc.Server
}

func New(addr string) api.Server {
	ctx := context.Background()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort("localhost", "9999"),
		opts...,
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	if err := build.RegisterGreetingHandler(ctx, mux, conn); err != nil {
		log.Fatalf("api/grpc/New: greeting handler registration has failed, %v", err)
	}

	// tlsCertificate, err := tls.X509KeyPair(cert.ServerCertificate, cert.ServerPrivateKey)
	// if err != nil {
	// 	log.Fatalf("api/grpc/New: tls cetificate setup has failed, %v", err)
	// }

	hs := http.Server{
		Addr:    addr,
		Handler: mux,
		// TLSConfig: &tls.Config{
		// 	Certificates: []tls.Certificate{tlsCertificate},
		// },
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	consoleLogger := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleLogger).With().Timestamp().Logger()

	grpcServer := grpcServer{
		logger:     &logger,
		restServer: &hs,
	}
	if err := grpcServer.startGRPCServer("localhost", "9999"); err != nil {
		grpcServer.logger.Error().Msgf("api/grpc/New: %v", err)
	}

	return &grpcServer
}

func (gs *grpcServer) ListenAndServe() error {
	return gs.restServer.ListenAndServe()
}

func (gs *grpcServer) Shutdown(ctx context.Context) error {
	err := gs.restServer.Shutdown(ctx)
	if err != nil {
		gs.logger.Err(err).Msgf("api/grpc/Shutdown: rest server shutdown error, %w", err)
	}

	grpcCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = gs.stopGRPCServer(grpcCtx)
	if err != nil {
		gs.logger.Err(err).Msgf("api/grpc/Shutdown: grpc server shutdown error, %w", err)
	}

	return err
}

func (gs *grpcServer) startGRPCServer(host, port string) error {
	gs.grpcServer = grpc.NewServer()
	build.RegisterGreetingServer(gs.grpcServer, greeting.New())

	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return fmt.Errorf("api/grpc/startGRPCServer: starting grpc server on host %s and port %s has failed: %w",
			host, port, err)
	}

	go func() {
		if err := gs.grpcServer.Serve(lis); err != nil {
			gs.logger.Error().Msgf("api/grpc/startGRPCServer: grpc server serve has failed: %v", err)
		}
	}()

	return nil
}

func (gs *grpcServer) stopGRPCServer(ctx context.Context) error {
	ok := make(chan struct{})

	go func() {
		gs.grpcServer.GracefulStop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		gs.grpcServer.Stop()
		return ctx.Err()
	}
}

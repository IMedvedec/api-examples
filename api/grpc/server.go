package grpc

import (
	"context"
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
	logger *zerolog.Logger
	server *http.Server
}

func New(addr string) api.Server {
	ctx := context.Background()

	gs := grpc.NewServer()
	build.RegisterGreetingServer(gs, greeting.New())

	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", "9999"))
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

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
		logger: &logger,
		server: &hs,
	}

	return &grpcServer
}

func (gs *grpcServer) ListenAndServe() error {
	return gs.server.ListenAndServe()
}

func (gs *grpcServer) Shutdown(ctx context.Context) error {
	return gs.server.Shutdown(ctx)
}

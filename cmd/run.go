package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/imedvedec/api-examples/api"
	"github.com/imedvedec/api-examples/api/builtin"
	"github.com/imedvedec/api-examples/api/grpc"
	"github.com/rs/zerolog"
)

// API server addresses.
const (
	builtinAPIaddr string = "localhost:8080"
	grpcAPIaddr    string = "localhost:8081"
)

// Server parametrisation.
const (
	shutdownDeadline time.Duration = 5 * time.Second
)

// application represents the main application with dependencies.
type application struct {
	logger *zerolog.Logger
}

// New is an application constructor.
func New() *application {
	consoleLogger := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleLogger).With().Timestamp().Logger()

	return &application{
		logger: &logger,
	}
}

// Run represents the application starting point.
func (app *application) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signals
		app.logger.Info().Msg(fmt.Sprintf("OS signal happened: %v", sig))
		cancel()
	}()

	app.serverLifeCycle(ctx)

	app.logger.Info().Msg("Application has finished successfully")
}

// serverLifeCycle is a helper function which manages API server
// setup, start and graceful shutdown.
func (app *application) serverLifeCycle(ctx context.Context) {
	servers := []api.Server{
		builtin.New(builtinAPIaddr),
		grpc.New(grpcAPIaddr),
	}

	for _, s := range servers {
		go func(s api.Server) {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("Error happened on server listen: %v", err))
			}

			//app.logger.Info().Msg(fmt.Sprintf("API started on: '%s'", builtinAPIaddr))
		}(s)
	}

	<-ctx.Done()

	app.logger.Info().Msg(fmt.Sprintf("Builtin API ('%s') has been stopped", builtinAPIaddr))

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownDeadline)
	defer cancel()

	shutdownWG := sync.WaitGroup{}
	for _, s := range servers {
		shutdownWG.Add(1)
		go func(s api.Server) {
			if err := s.Shutdown(ctxShutdown); err != nil {
				//app.logger.Panic().Msg(fmt.Sprintf("Error has happened on server shutdown: %v", builtinAPIaddr, err))
			}
			shutdownWG.Done()
		}(s)
	}
	shutdownWG.Wait()
}

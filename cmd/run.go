package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/imedvedec/api-examples/api/builtin"
)

// API server addresses.
const (
	builtinAPIaddr string = "localhost:8080"
)

// Server parametrisation.
const (
	shutdownDeadline time.Duration = 5 * time.Second
)

// Run represents the application starting point.
func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		sig := <-signals
		log.Printf("OS signal happened: %v\n", sig)
		cancel()
	}()

	serverLifeCycle(ctx)
	log.Println("Application finished")
}

// serverLifeCycle is a helper function which manages API server
// setup, start and graceful shutdown.
func serverLifeCycle(ctx context.Context) {
	server := builtin.New(builtinAPIaddr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Error happened on builtinAPI listen: %v", err))
		}
	}()

	log.Println(fmt.Sprintf("Builtin API started on: %s", builtinAPIaddr))

	<-ctx.Done()

	log.Println(fmt.Sprintf("Builtin API has been stopped"))

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownDeadline)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		panic(fmt.Sprintf("Error happened on bultinAPI shutdown: %v", err))
	}

}

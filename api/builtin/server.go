package builtin

import (
	"context"
	"net/http"

	"github.com/imedvedec/api-examples/api"
	"github.com/rs/zerolog"
)

// builtinServer represents a builtin server with its dependencies.
type builtinServer struct {
	logger *zerolog.Logger
	server *http.Server
}

// New is an builtin api server constructor.
func New(addr string) api.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/json", jsonHandler)
	mux.Handle("/middleware/json", clientJSONCheck(http.HandlerFunc(jsonHandler)))
	mux.Handle("/cookie/json", domainCookieSetHandler(http.HandlerFunc(jsonHandler)))
	mux.Handle("/cookie/json/secure", domainAndPathCookieSetHandler(http.HandlerFunc(jsonHandler)))

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	consoleLogger := zerolog.NewConsoleWriter()
	logger := zerolog.New(consoleLogger).With().Timestamp().Logger()

	builtinServer := builtinServer{
		logger: &logger,
		server: &server,
	}

	return &builtinServer
}

// ListenAndServe imeplements the api.Server.ListenAndServe method.
func (bs *builtinServer) ListenAndServe() error {
	return bs.server.ListenAndServe()
}

// Shutdown imeplements the api.Server.Shutdown method.
func (bs *builtinServer) Shutdown(ctx context.Context) error {
	return bs.server.Shutdown(ctx)
}

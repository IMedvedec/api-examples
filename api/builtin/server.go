package builtin

import (
	"context"
	"net/http"

	"github.com/imedvedec/api-examples/api"
	"github.com/rs/zerolog"
)

type builtinServer struct {
	logger *zerolog.Logger
	server *http.Server
}

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

func (bs *builtinServer) ListenAndServe() error {
	return bs.server.ListenAndServe()
}

func (bs *builtinServer) Shutdown(ctx context.Context) error {
	return bs.server.Shutdown(ctx)
}

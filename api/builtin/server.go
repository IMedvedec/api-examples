package builtin

import (
	"net/http"
)

func New(addr string) *http.Server {
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

	return &server
}

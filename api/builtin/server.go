package builtin

import (
	"net/http"
)

func New(addr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/json", jsonHandler)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &server
}

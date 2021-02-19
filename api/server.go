package api

import "context"

// Server defines a application server.
type Server interface {
	// ListenAndServe represents a http-like listen and server method.
	ListenAndServe() error
	// Shutdown is used server graceful shutdown.
	Shutdown(context.Context) error
}

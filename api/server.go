package api

import "context"

type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

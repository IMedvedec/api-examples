package cmd

import "github.com/imedvedec/api-examples/api/builtin"

const (
	builtinAPIaddr string = "localhost:8080"
)

func Run() {
	server := builtin.New(builtinAPIaddr)
	server.ListenAndServe()
}

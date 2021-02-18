package builtin

import (
	"fmt"
	"log"
	"net/http"
)

//preorderMiddleware is a preorder middleware example.
func preorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("preorder middleware execution!")
		next.ServeHTTP(w, r)
	})
}

//postorderMiddleware is a postorder middleware example.
func postorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("postorder middleware execution!")
	})
}

//clientJSONCheck checks that the client can accept JSON responses.
func clientJSONCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(accept)

		if headerValue == "" {
			fmt.Fprintf(w, "Accept header is not set!")
			return
		}

		if headerValue != contentTypeJSON {
			fmt.Fprintf(w, "Sorry, this endpoint is JSON only!")
			return
		}

		next.ServeHTTP(w, r)
	})
}

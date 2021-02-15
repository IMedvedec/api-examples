package builtin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(
		w,
		r,
		[]byte("Hello, this is a simple handler"),
		http.StatusOK,
		map[string]string{
			contentType: contentTypeTextHTML,
		},
	)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	type exampleInfo struct {
		ID   int64  `json:"id"`
		Data string `json:"data"`
	}

	jsonString, err := json.Marshal(
		exampleInfo{
			ID:   1,
			Data: "My data example",
		},
	)
	if err != nil {
		writeResponse(
			w,
			r,
			[]byte(fmt.Sprintf("Internal server error: %s", err)),
			http.StatusInternalServerError,
			map[string]string{
				contentType: contentTypeTextHTML,
			},
		)
		return
	}

	writeResponse(
		w,
		r,
		jsonString,
		http.StatusOK,
		map[string]string{
			contentType: contentTypeJSON,
		},
	)
}

func writeResponse(
	w http.ResponseWriter,
	r *http.Request,
	message []byte,
	status int,
	headers map[string]string,
) {
	for k, v := range headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(status)
	w.Write([]byte(message))
}

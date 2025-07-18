package pltt

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
)

func StartServer() {
	http.HandleFunc("POST /{$}", writeReturnedError(keysPost))
	http.HandleFunc("GET /{keys}", writeReturnedError(dataGet))
	http.HandleFunc("POST /{key}", writeReturnedError(dataPost))
	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		slog.Error("Unable to start service", "error", err)
	}
}

func writeReturnedError(handler func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			writeError(w, err)
		}
	}
}

func writeError(w http.ResponseWriter, err error) {
	is := func(target error) bool {
		return errors.Is(err, target)
	}
	code := 0
	switch {
	case is(errorInvalidBody) || is(errorInvalidFormat) ||
		is(errorInvalidKey) || is(errorInvalidKeyCount):
		code = 400
	case is(errorKeyNotFound):
		code = 404
	default:
		code = 500
		slog.Error("Server error", "error", err)
	}
	w.WriteHeader(code)
	io.WriteString(w, err.Error())
}

package pltt

import (
	"io"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	StatusCode int
	Text       string
}

func (e errorResponse) Error() string {
	return e.Text
}

var (
	invalidPathError = errorResponse{400, "Invalid path. Only alphanumeric characters, underscores & dashes are supported"}
	invalidBodyError = errorResponse{400, "Invalid body. Only numbers are supported"}
)

func writeError(w http.ResponseWriter, e error) {
	er, iser := e.(errorResponse)
	if iser {
		w.WriteHeader(er.StatusCode)
	} else {
		w.WriteHeader(500)
		slog.Error("Internal server error", "error", e)
	}
	io.WriteString(w, e.Error())
}

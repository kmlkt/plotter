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
	invalidKeyError    = errorResponse{400, "Invalid key. Only alphanumeric characters, underscores & dashes are supported"}
	invalidFormatError = errorResponse{400, "Invalid response format. Supported values: .html (default), .svg, .csv and .json"}
	invalidBodyError   = errorResponse{400, "Invalid body. Only numbers are supported"}
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

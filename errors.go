package pltt

import (
	"fmt"
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
	errorInvalidKey      = errorResponse{400, "Invalid key. Only alphanumeric characters, underscores & dashes are supported"}
	errorInvalidFormat   = errorResponse{400, "Invalid response format. Supported values: .html (default), .svg, .csv and .json"}
	errorInvalidBody     = errorResponse{400, "Invalid body. Only numbers are supported"}
	errorInvalidKeyCount = errorResponse{400, "Only 1 key is supported in POST request"}
)

func errorCantAccess(hash string, hashMode byte, requestMode byte) error {
	return errorResponse{
		403,
		fmt.Sprintf("Key hash %s can only %c, but you're trying to %c", hash, hashMode, requestMode),
	}
}

func stringifyMode(mode byte) string {
	if mode == 'r' {
		return "read"
	} else {
		return "write"
	}
}

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

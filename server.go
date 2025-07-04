package pltt

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func StartServer() {
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		err := post(r)
		if err != nil {
			writeError(w, err)
		}
	})
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ans, err := get(r)
		if err != nil {
			writeError(w, err)
			return
		}
		err = json.NewEncoder(w).Encode(ans)
		if err != nil {
			slog.Error("Cannot marshal json", "error", err)
		}
	})
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		slog.Error("Unable to start service", "error", err)
	}
}

var keyValidator = regexp.MustCompile("^[A-Za-z0-9-]+$")

func post(r *http.Request) error {
	key := strings.Trim(r.URL.EscapedPath(), "/")
	if !keyValidator.MatchString(key) {
		return invalidPathError
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	value, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return invalidBodyError
	}
	return write(key, value)
}

func get(r *http.Request) (any, error) {
	key := strings.Trim(r.URL.EscapedPath(), "/")
	if !keyValidator.MatchString(key) {
		return nil, invalidPathError
	}
	iter, err := read(key)
	if err != nil {
		return nil, err
	}
	ans := iter
	if r.URL.Query().Has("since") {
		t, err := time.Parse(time.DateTime, r.URL.Query().Get("since"))
		if err != nil {
			return nil, err
		}
		ans = since(ans, t)
	}
	return slices.Collect(ans), nil
}

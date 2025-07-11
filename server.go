package pltt

import (
	"fmt"
	"io"
	"iter"
	"log/slog"
	"net/http"
	"regexp"
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
		err := get(w, r)
		if err != nil {
			writeError(w, err)
			return
		}
	})
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		slog.Error("Unable to start service", "error", err)
	}
}

func post(r *http.Request) error {
	keys, _, err := validateUrl(r)
	if err != nil {
		return err
	}
	if len(keys) != 1 {
		return errorInvalidKeyCount
	}
	key := keys[0]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	value, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return errorInvalidBody
	}
	return write(key, value)
}

func get(w http.ResponseWriter, r *http.Request) error {
	keys, format, err := validateUrl(r)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType(format))
	data := make([]iter.Seq[record], len(keys))
	for i, key := range keys {
		data[i], err = read(key)
		if err != nil {
			return err
		}
	}

	err = applyQueryFilters(r, data)
	if err != nil {
		return err
	}
	marshalData(format, data, w, keys)

	return nil
}

var keyValidator = regexp.MustCompile("^[A-Za-z0-9_-]+$")

type format string

const (
	formatHtml format = "html"
	formatSvg  format = "svg"
	formatCsv  format = "csv"
	formatJson format = "json"
)

func validateUrl(r *http.Request) ([]string, format, error) {
	keys, f := parseUrl(r)
	if len(keys) == 0 {
		return nil, "", errorInvalidKey
	}
	if f != formatHtml && f != formatSvg && f != formatCsv && f != formatJson {
		return nil, "", errorInvalidFormat
	}
	for i, key := range keys {
		key, err := parseKey(key)
		if err != nil {
			return nil, "", err
		}
		keys[i] = key
	}
	return keys, f, nil
}

func parseKey(key string) (string, error) {
	if len(key) == 0 {
		return "", errorInvalidKey
	}
	if !keyValidator.MatchString(key) {
		return "", errorInvalidKey
	}
	return key, nil
}

func parseUrl(r *http.Request) ([]string, format) {
	parts := strings.Split(strings.Trim(r.URL.EscapedPath(), "/"), ".")
	switch len(parts) {
	case 0:
		return nil, ""
	case 1:
		return strings.Split(parts[0], "&"), "html"
	case 2:
		return strings.Split(parts[0], "&"), format(parts[1])
	default:
		return nil, ""
	}
}

func contentType(f format) string {
	switch f {
	case formatHtml:
		return "text/html"
	case formatSvg:
		return "image/svg"
	case formatCsv:
		return "text/csv"
	case formatJson:
		return "application/json"
	default:
		return ""
	}
}

func applyQueryFilters(r *http.Request, data []iter.Seq[record]) error {
	if r.URL.Query().Has("since") {
		t, err := time.Parse("2006-01-02T15:04:05", r.URL.Query().Get("since"))
		if err != nil {
			return err
		}
		for i, _ := range data {
			data[i] = since(data[i], t)
		}
	}
	return nil
}

func marshalData(format format, data []iter.Seq[record], w io.Writer, titles []string) {
	switch format {
	case formatHtml:
		fmt.Fprintf(w, `
			<!doctype html>
			<html>
			<head>
				<title>%s - plotter</title>
				<style>
					html, body {
						margin: 0;
						padding: 0;
						overflow: clip;
					}
					svg {
						width: 90vw;
						height: 90vh;
						margin: 10vh 0 0 0;
						border: 0;
						padding: 0;
					}
				</style>
			</head>
			<body>
			`, strings.Join(titles, " & "))
		buildGraph(newRecordsDescriptor(data), titles, w)
		fmt.Fprint(w, `
			</body>
			</html>
			`)
	case formatSvg:
		buildGraph(newRecordsDescriptor(data), titles, w)
	case formatCsv:
		fmt.Fprintf(w, "Time,Value\n")
		for _, ds := range data {
			for r := range ds {
				r.Fprintf(w, "%s,%f\n")
			}
		}
	case formatJson:
		fmt.Fprint(w, "{")
		first := true
		for _, ds := range data {
			for r := range ds {
				if first {
					first = false
				} else {
					fmt.Fprint(w, ",")
				}
				r.Fprintf(w, `"%s": %f`)
			}
		}
		fmt.Fprint(w, "}")
	}
}

func (r record) Fprintf(w io.Writer, format string) {
	fmt.Fprintf(w, format, r.UTCTimeString(), r.Value)
}

func (r record) UTCTimeString() string {
	return r.Timestamp.UTC().Format(time.DateTime)
}

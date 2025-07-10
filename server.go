package pltt

import (
	"fmt"
	"io"
	"iter"
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
	key := strings.Trim(r.URL.EscapedPath(), "/")
	if !keyValidator.MatchString(key) {
		return invalidKeyError
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

func get(w http.ResponseWriter, r *http.Request) error {
	key, format, err := validateUrl(r)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType(format))

	data, err := read(key)
	if err != nil {
		return err
	}

	err = applyQueryFilters(r, &data)
	if err != nil {
		return err
	}
	marshalData(format, data, w, key)

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

func validateUrl(r *http.Request) (string, format, error) {
	key, f := parseUrl(r)
	if !keyValidator.MatchString(key) {
		return "", "", invalidKeyError
	}
	if f != formatHtml && f != formatSvg && f != formatCsv && f != formatJson {
		return "", "", invalidFormatError
	}
	return key, f, nil
}

func parseUrl(r *http.Request) (string, format) {
	parts := strings.Split(strings.Trim(r.URL.EscapedPath(), "/"), ".")
	switch len(parts) {
	case 0:
		return "WTF?!", ""
	case 1:
		return parts[0], "html"
	case 2:
		return parts[0], format(parts[1])
	default:
		return "WTF?!", ""
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

func applyQueryFilters(r *http.Request, data *iter.Seq[record]) error {
	if r.URL.Query().Has("since") {
		t, err := time.Parse("2006-01-02T15:04:05", r.URL.Query().Get("since"))
		if err != nil {
			return err
		}
		*data = since(*data, t)
	}
	return nil
}

func marshalData(format format, data iter.Seq[record], w io.Writer, title string) {
	switch format {
	case formatHtml:
		fmt.Fprintf(w, `
			<!doctype html>
			<html>
			<head>
				<title>%s @ plotter</title>
				<style>
					html {
						margin: 10vh 0 0 0;
					}
					svg {
						width: 90vw;
						height: 90vh;
					}
				</style>
			</head>
			<body>
			`, title)
		buildGraph(newRecordsDescriptor(slices.Collect(data)), w)
		fmt.Fprint(w, `
			</body>
			</html>
			`)
	case formatSvg:
		buildGraph(newRecordsDescriptor(slices.Collect(data)), w)
	case formatCsv:
		fmt.Fprintf(w, "Time,Value\n")
		for r := range data {
			r.Fprintf(w, "%s,%f\n")
		}
	case formatJson:
		fmt.Fprint(w, "{")
		first := true
		for r := range data {
			if first {
				first = false
			} else {
				fmt.Fprint(w, ",")
			}
			r.Fprintf(w, `"%s": %f`)
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

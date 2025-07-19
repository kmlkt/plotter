package pltt

import (
	"net/http"
	"regexp"
	"strings"
	"time"
)

type graphConfig struct {
	format format
	keys   []string
	minT   time.Time
	maxT   time.Time
	sumD   time.Duration
}

var keyValidator = regexp.MustCompile("^[A-Za-z0-9_-]+$")

type format string

const (
	formatHtml format = "html"
	formatSvg  format = "svg"
	formatCsv  format = "csv"
	formatJson format = "json"
)

func validateUrl(r *http.Request, cfg *graphConfig) error {
	cfg.keys, cfg.format = parseUrl(r)
	if len(cfg.keys) == 0 {
		return errorInvalidKey
	}
	if cfg.format != formatHtml && cfg.format != formatSvg && cfg.format != formatCsv && cfg.format != formatJson {
		return errorInvalidFormat
	}
	for _, key := range cfg.keys {
		err := validateKey(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseUrl(r *http.Request) ([]string, format) {
	pathKeys := r.PathValue("keys")
	if pathKeys == "" {
		pathKeys = r.PathValue("key")
	}
	parts := strings.Split(pathKeys, ".")
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

func validateKey(key string) error {
	if len(key) == 0 {
		return errorInvalidKey
	}
	if !keyValidator.MatchString(key) {
		return errorInvalidKey
	}
	return nil
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

func parseQueryArgs(r *http.Request, cfg *graphConfig) error {
	query := r.URL.Query()
	cfg.minT = time.Unix(0, 0)
	cfg.maxT = time.Now()
	cfg.sumD = time.Duration(0)
	var err error
	if query.Has("since") {
		cfg.minT, err = time.Parse("2006-01-02T15:04:05", query.Get("since"))
		if err != nil {
			return err
		}
	}
	if query.Has("until") {
		cfg.maxT, err = time.Parse("2006-01-02T15:04:05", query.Get("until"))
		if err != nil {
			return err
		}
	}
	if query.Has("last") {
		d, err := time.ParseDuration(query.Get("last"))
		if err != nil {
			return err
		}
		cfg.minT = cfg.maxT.Add(-d)
	}
	if query.Has("sum") {
		cfg.sumD, err = time.ParseDuration(query.Get("sum"))
		if err != nil {
			return err
		}
	}
	return nil
}

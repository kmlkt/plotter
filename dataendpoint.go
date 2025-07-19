package plotter

import (
	_ "embed"
	"fmt"
	"io"
	"iter"
	"net/http"
	"strings"
	"time"
)

func dataGet(w http.ResponseWriter, r *http.Request) error {
	cfg := graphConfig{}
	cfg.minT = time.Unix(0, 0)
	cfg.maxT = time.Now().UTC()
	cfg.sumD = time.Duration(0)
	cfg.xLabels = 10
	cfg.yLabels = 10
	err := validateUrl(r, &cfg)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", contentType(cfg.format))
	err = parseQueryArgs(r, &cfg)
	if err != nil {
		return err
	}
	data := make([]iter.Seq[record], len(cfg.keys))

	for i, key := range cfg.keys {
		data[i], err = read(key)
		if err != nil {
			return err
		}
	}
	applyQueryFilters(data, cfg)
	marshalData(w, data, cfg)
	return nil
}

func dataPost(_ http.ResponseWriter, r *http.Request) error {
	cfg := graphConfig{}
	err := validateUrl(r, &cfg)
	if err != nil {
		return err
	}
	if len(cfg.keys) != 1 {
		return formatError(errorInvalidKeyCount, len(cfg.keys))
	}
	key := cfg.keys[0]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	value, err := parseDecimal(string(body))
	if err != nil {
		return formatError(errorInvalidBody, string(body), err.Error())
	}
	return write(key, value)
}

//go:embed index.html
var indexHtml string

func marshalData(w io.Writer, data []iter.Seq[record], cfg graphConfig) {
	switch cfg.format {
	case formatHtml:
		fmt.Fprintf(w, indexHtml, strings.Join(cfg.keys, " & "))
	case formatSvg:
		buildGraph(w, data, cfg)
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

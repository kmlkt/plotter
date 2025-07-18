package pltt

import (
	"fmt"
	"io"
	"iter"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func dataGet(w http.ResponseWriter, r *http.Request) error {
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
	//fmt.Errorf("%w", a ...any)
	err = applyQueryFilters(r, data)
	if err != nil {
		return err
	}
	marshalData(format, data, w, keys)

	return nil
}

func dataPost(_ http.ResponseWriter, r *http.Request) error {
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

func parseKey(key string) (string, error) {
	if len(key) == 0 {
		return "", errorInvalidKey
	}
	if !keyValidator.MatchString(key) {
		return "", errorInvalidKey
	}
	return key, nil
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
	filters := make([]func(iter.Seq[record]) iter.Seq[record], 0)
	query := r.URL.Query()
	if query.Has("since") {
		t, err := time.Parse("2006-01-02T15:04:05", query.Get("since"))
		if err != nil {
			return err
		}
		filters = append(filters, func(s iter.Seq[record]) iter.Seq[record] {
			return since(s, t)
		})
	}
	if query.Has("until") {
		t, err := time.Parse("2006-01-02T15:04:05", query.Get("until"))
		if err != nil {
			return err
		}
		filters = append(filters, func(s iter.Seq[record]) iter.Seq[record] {
			return until(s, t)
		})
	}
	if query.Has("last") {
		d, err := time.ParseDuration(query.Get("last"))
		if err != nil {
			return err
		}
		filters = append(filters, func(s iter.Seq[record]) iter.Seq[record] {
			return last(s, d)
		})
	}
	if query.Has("sum") {
		d, err := time.ParseDuration(query.Get("sum"))
		if err != nil {
			return err
		}
		filters = append(filters, func(s iter.Seq[record]) iter.Seq[record] {
			return intervalSum(s, d)
		})
	}

	for i, _ := range data {
		for _, f := range filters {
			data[i] = f(data[i])
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
					#tooltip {
						position: absolute;
					}
				</style>
			</head>
			<body>
			`, strings.Join(titles, " & "))
		buildGraph(newRecordsDescriptor(data), titles, w)
		fmt.Fprint(w, `
			<span id='tooltip'></span>
			<script>
			const $svg = document.querySelector('svg')
			const $tooltip = document.querySelector('#tooltip')
			const circles = Array.from(document.querySelectorAll('circle')).map(c => ({
				x: (c.getBoundingClientRect().left + c.getBoundingClientRect().right) / 2,
				y: c.getBoundingClientRect().y,
				text: c.textContent,
			}))
			$svg.addEventListener('mousemove', e => {
				const x = e.clientX
				const y = e.clientY
				const dist = (c) => Math.abs(x - c.x)
				let nearest = undefined
				for (const c of circles) {
					if (!nearest
					|| dist(c) < dist(nearest)
					|| dist(c) == dist(nearest) && c.y < nearest.y) {
						nearest = c
					}
				}
				if(nearest){
					tooltip.textContent = nearest.text
					tooltip.style.left = (nearest.x - tooltip.clientWidth / 2) + 'px'
					tooltip.style.top = (nearest.y - 20) + 'px'
				}
			})
			</script>
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

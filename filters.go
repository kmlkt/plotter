package plotter

import (
	"iter"
	"time"
)

func applyQueryFilters(data []iter.Seq[record], cfg graphConfig) {
	for i, _ := range data {
		data[i] = intervalSum(until(since(data[i], cfg.minT), cfg.maxT), cfg)
	}
}

func since(s iter.Seq[record], t time.Time) iter.Seq[record] {
	return func(yield func(record) bool) {
		for r := range s {
			if !yield(r) {
				return
			}
			if r.Timestamp.Before(t) {
				return
			}
		}
	}
}

func until(s iter.Seq[record], t time.Time) iter.Seq[record] {
	return func(yield func(record) bool) {
		for r := range s {
			if !yield(r) {
				return
			}
			if r.Timestamp.After(t) {
				continue
			}
		}
	}
}

func intervalSum(s iter.Seq[record], cfg graphConfig) iter.Seq[record] {
	if cfg.sumD == 0 {
		return s
	}
	ans := make(map[time.Time]decimal)
	for r := range s {
		ans[r.Timestamp.Truncate(cfg.sumD)] += r.Value
	}
	return func(yield func(record) bool) {
		for t := cfg.minT.Truncate(cfg.sumD); t.Before(cfg.maxT); t = t.Add(cfg.sumD) {
			if !yield(record{t, ans[t]}) {
				return
			}
			if !yield(record{t.Add(cfg.sumD), ans[t]}) {
				return
			}
		}
	}
}

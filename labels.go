package plotter

import (
	"time"
)

type label struct {
	Value    string
	Position float64
}

// seconds
var timeIntervals = []time.Duration{time.Second, 10 * time.Second, 15 * time.Second, 30 * time.Second, time.Minute, 2 * time.Minute, 5 * time.Minute, 10 * time.Minute, 15 * time.Minute, 30 * time.Minute, time.Hour, 4 * time.Hour, 6 * time.Hour, 12 * time.Hour, 24 * time.Hour, 10 * 24 * time.Hour, 20 * 24 * time.Hour, 50 * 24 * time.Hour, 100 * 24 * time.Hour, 200 * 24 * time.Hour, 500 * 24 * time.Hour, 1000 * 24 * time.Hour}
var valueIntervals = []decimal{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000, 100000, 200000, 500000, 1000000, 2000000, 5000000, 10000000, 20000000, 50000000, 100000000, 200000000, 500000000, 1000000000, 2000000000, 5000000000, 10000000000, 20000000000, 50000000000, 100000000000, 200000000000, 500000000000, 1000000000000, 2000000000000, 5000000000000, 10000000000000, 20000000000000, 50000000000000, 100000000000000, 200000000000000, 500000000000000, 1000000000000000, 2000000000000000, 5000000000000000, 10000000000000000, 20000000000000000, 50000000000000000, 100000000000000000, 200000000000000000, 500000000000000000, 1000000000000000000, 2000000000000000000, 5000000000000000000}

func labels(data recordsDescriptor, cfg graphConfig) (times []label, values []label) {
	times = optimalLabels(data.Times, cfg.xLabels, timeIntervals, timeFormat)
	values = optimalLabels(data.Values, cfg.yLabels, valueIntervals, func(decimal) string { return "" })
	return
}

func optimalLabels[T calculable[T, D], D diff](values descriptor[T, D], maxCount int,
	intervals []D, format func(interval D) string) []label {
	interval, count := optimalInterval(values, int64(maxCount), intervals)
	ans := make([]label, count)
	v := firstLabel(values, interval)
	for i := range count {
		ans[i].Value = v.Format(format(interval))
		ans[i].Position = values.position(v)
		v = v.Add(interval)
	}
	return ans
}

func optimalInterval[T calculable[T, D], D diff](values descriptor[T, D], maxCount int64, intervals []D) (interval D, count int64) {
	for _, intr := range intervals {
		cnt := labelCount(values, intr)
		if count < cnt && cnt <= maxCount {
			interval = intr
			count = cnt
		}
	}
	return
}

func labelCount[T calculable[T, D], D diff](values descriptor[T, D], interval D) int64 {
	return int64(lastLabel(values, interval).Sub(firstLabel(values, interval)))/int64(interval) + 1
}

func firstLabel[T calculable[T, D], D diff](values descriptor[T, D], interval D) T {
	return values.Min.Add(interval - 1).Truncate(interval)
}

func lastLabel[T calculable[T, D], D diff](values descriptor[T, D], interval D) T {
	return values.Max.Truncate(interval)
}

func timeFormat(interval time.Duration) string {
	switch {
	case interval >= 24*time.Hour:
		return time.DateOnly
	case interval >= 6*time.Hour:
		return "01-02 15:04"
	case interval >= time.Minute:
		return "15:04"
	default:
		return time.TimeOnly
	}
}

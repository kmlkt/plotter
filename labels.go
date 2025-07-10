package pltt

import (
	"math"
	"strconv"
	"time"
)

type label struct {
	Value    string
	Position float64
}

// seconds
var timeIntervals = []float64{1, 10, 15, 30, 60, 120, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400, 864000}
var valueIntervals = []float64{1e-10, 2e-10, 5e-10, 1e-9, 2e-9, 5e-9, 1e-8, 2e-8, 5e-8, 1e-7, 2e-7, 5e-7, 1e-6, 2e-6, 5e-6, 1e-5, 2e-5, 5e-5, 1e-4, 2e-4, 5e-4, 1e-3, 2e-3, 5e-3, 1e-2, 2e-2, 5e-2, 1e-1, 2e-1, 5e-1, 1e0, 2e0, 5e0, 1e1, 2e1, 5e1, 1e2, 2e2, 5e2, 1e3, 2e3, 5e3, 1e4, 2e4, 5e4, 1e5, 2e5, 5e5, 1e6, 2e6, 5e6, 1e7, 2e7, 5e7, 1e8, 2e8, 5e8, 1e9, 2e9, 5e9, 1e10, 2e10, 5e10}

func labels(data recordsDescriptor) (times []label, values []label) {
	times = optimalLabels(data.Times, timeIntervals, TimeString)
	values = optimalLabels(data.Values, valueIntervals, ValueString)
	return
}

func optimalLabels(values descriptor[float64], intervals []float64, stringify func(interval float64, label float64) string) []label {
	interval, count := optimalInterval(values, intervals)
	ans := make([]label, count)
	for i := range count {
		v := interval * (math.Ceil(values.Min/interval) + float64(i))
		ans[i].Value = stringify(interval, v)
		ans[i].Position = (v - values.Min) / (values.Max - values.Min)
	}
	return ans
}

func optimalInterval(values descriptor[float64], intervals []float64) (interval float64, count int) {
	for _, intr := range intervals {
		cnt := labelCount(values, intr)
		if count < cnt && cnt <= 10 {
			interval = intr
			count = cnt
		}
	}
	return
}

func labelCount(values descriptor[float64], interval float64) int {
	firstIndex := int(math.Ceil(values.Min / interval))
	lastIndex := int(values.Max / interval)
	return lastIndex - firstIndex + 1
}

func TimeString(interval float64, label float64) string {
	layout := time.TimeOnly
	if interval >= 60 {
		layout = "15:04"
	}
	if interval >= 6*60*60 {
		layout = time.DateTime
	}
	if interval >= 24*60*60 {
		layout = time.DateOnly
	}
	time := time.Unix(int64(label), 0)
	return time.UTC().Format(layout)
}

func ValueString(_ float64, label float64) string {
	return strconv.FormatFloat(label, 'G', -1, 64)
}

package pltt

import (
	"slices"
	"time"
)

type label struct {
	Value    string
	Position float64
}

// seconds
var timeIntervals = []int64{1, 10, 15, 30, 60, 120, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400, 864000}
var valueIntervals = []int64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000, 100000, 200000, 500000, 1000000, 2000000, 5000000, 10000000, 20000000, 50000000, 100000000, 200000000, 500000000, 1000000000, 2000000000, 5000000000, 10000000000, 20000000000, 50000000000, 100000000000, 200000000000, 500000000000, 1000000000000, 2000000000000, 5000000000000, 10000000000000, 20000000000000, 50000000000000, 100000000000000, 200000000000000, 500000000000000, 1000000000000000, 2000000000000000, 5000000000000000, 10000000000000000, 20000000000000000, 50000000000000000, 100000000000000000, 200000000000000000, 500000000000000000, 1000000000000000000, 2000000000000000000, 5000000000000000000}

func labels(data recordsDescriptor) (times []label, values []label) {
	times = optimalLabels(data.Times, timeIntervals, timeString)
	values = optimalLabels(data.Values, valueIntervals, valueString)
	return
}

func optimalLabels(values descriptor[int64], intervals []int64, stringify func(interval int64, label int64) string) []label {
	interval, count := optimalInterval(values, intervals)
	ans := make([]label, count)
	v := values.Min + interval - 1
	v -= v % interval
	for i := range count {
		ans[i].Value = stringify(interval, v)
		ans[i].Position = float64(v-values.Min) / float64(values.Max-values.Min)
		v += interval
	}
	return ans
}

func optimalInterval(values descriptor[int64], intervals []int64) (interval int64, count int) {
	for _, intr := range intervals {
		cnt := labelCount(values, intr)
		if int64(count) < cnt && cnt <= 10 {
			interval = intr
			count = int(cnt)
		}
	}
	return
}

func labelCount(values descriptor[int64], interval int64) int64 {
	firstIndex := (values.Min + interval - 1) / interval
	lastIndex := values.Max / interval
	return lastIndex - firstIndex + 1
}

func timeString(interval int64, label int64) string {
	layout := time.TimeOnly
	if interval >= 60 {
		layout = "15:04"
	}
	if interval >= 6*60*60 {
		layout = "01-02 15:04"
	}
	if interval >= 24*60*60 {
		layout = time.DateOnly
	}
	time := time.Unix(int64(label), 0)
	return time.UTC().Format(layout)
}

func valueString(_ int64, label int64) string {
	neg := label < 0
	if neg {
		label = -label
	}
	chars := make([]byte, 0)
	writeDigit := func() {
		chars = append(chars, '0'+byte(label%10))
		label /= 10
	}
	for range 9 {
		writeDigit()
	}
	chars = append(chars, '.')
	writeDigit()
	for label != 0 {
		writeDigit()
	}

	for chars[0] == '0' {
		chars = chars[1:]
	}
	if chars[0] == '.' {
		chars = chars[1:]
	}
	if neg {
		chars = append(chars, '-')
	}
	slices.Reverse(chars)
	return string(chars)
}

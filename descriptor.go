package pltt

import (
	"cmp"
	"iter"
	"math"
	"slices"
)

type descriptor[T cmp.Ordered] struct {
	Min T
	Max T
}

func (d *descriptor[T]) Append(value T) {
	d.Max = max(d.Max, value)
	d.Min = min(d.Min, value)
}

type recordsDescriptor struct {
	Data [][]record
	// seconds
	Times descriptor[int64]
	// x1e9
	Values descriptor[int64]
}

var emptyDescriptor = descriptor[int64]{math.MaxInt64, math.MinInt64}

const e9 = 1_000_000_000

func newRecordsDescriptor(iters []iter.Seq[record]) recordsDescriptor {
	data := make([][]record, len(iters))
	t := emptyDescriptor
	v := emptyDescriptor
	for i, iter := range iters {
		data[i] = slices.Collect(iter)
		for _, r := range data[i] {
			t.Append(r.Timestamp.Unix())
			v.Append(describeValue(r.Value))
		}
	}
	if v.Min == math.MaxInt64 {
		v.Min = 0
		v.Max = e9
	}
	if t.Min == math.MaxInt64 {
		t.Min = 0
		t.Max = 1
	}
	if v.Min == v.Max {
		v.Min -= e9
		v.Max += e9
	}
	if t.Min == t.Max {
		t.Min -= 1
		t.Max += 1
	}
	return recordsDescriptor{data, t, v}
}

func describeValue(value float64) int64 {
	return int64(value * 1e9)
}

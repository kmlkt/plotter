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
	Times  descriptor[float64]
	Values descriptor[float64]
}

var emptyDescriptor = descriptor[float64]{math.MaxFloat64, -math.MaxFloat64}

func newRecordsDescriptor(iters []iter.Seq[record]) recordsDescriptor {
	data := make([][]record, len(iters))
	t := emptyDescriptor
	v := emptyDescriptor
	for i, iter := range iters {
		data[i] = slices.Collect(iter)
		for _, r := range data[i] {
			t.Append(r.TimeFloat64())
			v.Append(r.Value)
		}
	}
	return recordsDescriptor{data, t, v}
}

func (r record) TimeFloat64() float64 {
	return float64(r.Timestamp.UnixNano()) / 1e9
}

package pltt

import (
	"cmp"
	"slices"
)

type descriptor[T any] struct {
	Min T
	Max T
}

func newDescriptor[T cmp.Ordered](values []T) descriptor[T] {
	if len(values) == 0 {
		return descriptor[T]{}
	}
	return descriptor[T]{slices.Min(values), slices.Max(values)}
}

type recordsDescriptor struct {
	Data []record
	// seconds
	Times  descriptor[float64]
	Values descriptor[float64]
}

func newRecordsDescriptor(data []record) recordsDescriptor {
	times := make([]float64, len(data))
	values := make([]float64, len(data))
	for i, r := range data {
		times[i] = r.TimeFloat64()
		values[i] = r.Value
	}
	return recordsDescriptor{data, newDescriptor(times), newDescriptor(values)}
}

func (r record) TimeFloat64() float64 {
	return float64(r.Timestamp.UnixNano()) / 1e9
}

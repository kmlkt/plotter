package plotter

import "time"

type diff interface{ ~int64 }

type calculable[Self any, D diff] interface {
	Add(m D) Self
	Sub(u Self) D
	Truncate(m D) Self
	String() string
	Format(format string) string
}

type descriptor[T calculable[T, D], D diff] struct {
	Min T
	Max T
}

func (d descriptor[T, D]) position(value T) float64 {
	return float64(value.Sub(d.Min)) / float64(d.Max.Sub(d.Min))
}

func newDecriptor[T calculable[T, D], D diff](v T) *descriptor[T, D] {
	return &descriptor[T, D]{v, v}
}

func (d *descriptor[T, D]) Append(value T) {
	if value.Sub(d.Max) > 0 {
		d.Max = value
	}
	if value.Sub(d.Min) < 0 {
		d.Min = value
	}
}

type recordsDescriptor struct {
	Times  descriptor[time.Time, time.Duration]
	Values descriptor[decimal, decimal]
}

func newRecordsDescriptor(data [][]record) recordsDescriptor {
	var t *descriptor[time.Time, time.Duration]
	var v *descriptor[decimal, decimal]
	for _, di := range data {
		for _, r := range di {
			if t == nil || v == nil {
				t = newDecriptor(r.Timestamp)
				v = newDecriptor(r.Value)
			} else {
				t.Append(r.Timestamp)
				v.Append(r.Value)
			}
		}
	}
	if t == nil || v == nil {
		t = newDecriptor(time.Unix(0, 0))
		v = newDecriptor(decimal(0))
	}
	if v.Min == v.Max {
		v.Min -= d1
		v.Max += d1
	}
	if t.Min.Equal(t.Max) {
		t.Min.Add(-time.Second)
		t.Max.Add(+time.Second)
	}
	return recordsDescriptor{*t, *v}
}

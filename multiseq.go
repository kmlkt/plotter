package pltt

import (
	"iter"
	"slices"
)

func multiCollect[T any](iters []iter.Seq[T]) [][]T {
	ans := make([][]T, len(iters))
	for i, it := range iters {
		ans[i] = slices.Collect(it)
	}
	return ans
}

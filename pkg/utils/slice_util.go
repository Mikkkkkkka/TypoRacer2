package utils

import "golang.org/x/exp/constraints"

func Sum[T constraints.Integer | constraints.Float](arr []T) float64 {
	var sum float64
	for _, v := range arr {
		sum += float64(v)
	}
	return sum
}

func Average[T constraints.Integer | constraints.Float](arr []T) float64 {
	return Sum(arr) / float64(len(arr))
}

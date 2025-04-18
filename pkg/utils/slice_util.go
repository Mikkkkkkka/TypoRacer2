package utils

import "golang.org/x/exp/constraints"

func Sum[T constraints.Integer](arr []T) int64 {
	var sum int64
	for _, v := range arr {
		sum += int64(v)
	}
	return sum
}

func Average[T constraints.Integer](arr []T) float64 {
	return float64(Sum(arr)) / float64(len(arr))
}

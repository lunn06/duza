package bit

import "golang.org/x/exp/constraints"

func ByIndex[T constraints.Integer](b T, idx int) T {
	return (b >> idx) & 1
}

func ClearFirstN[T constraints.Integer](b T, n int) T {
	return b >> n << n
}

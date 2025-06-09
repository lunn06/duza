package stego

import (
	"math/bits"
)

const (
	Delta           = 1 << 0
	InfoPixelsCount = bits.UintSize / Delta
)

func GetBit[T rune | byte | uint | uint16](b T, idx int) T {
	return (b >> idx) & 0x01
}

func ClearLastBits[T byte | uint32 | uint16](b T, n int) T {
	return b >> n << n
}

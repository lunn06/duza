package stego

const (
	Delta     = 1 << 0
	OffsetLen = 64
	UTF8Len   = 3
)

func GetBit[T int | rune | byte | uint | uint16](b T, idx int) T {
	return (b >> idx) & 0x01
}

func ClearLastBits[T byte | uint32 | uint16](b T, n int) T {
	return b >> n << n
}

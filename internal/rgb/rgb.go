package rgb

import "image/color"

type RGB int

const (
	Undefined RGB = iota
	Red
	Green
	Blue
)

func (c RGB) NRGBA64Value(color color.NRGBA64) uint16 {
	switch c {
	case Red:
		return color.R
	case Green:
		return color.G
	case Blue:
		return color.B
	case Undefined:
		fallthrough
	default:
		panic("unexpected color")
	}
}

func (c RGB) Next() RGB {
	switch c {
	case Red:
		return Green
	case Green:
		return Blue
	case Blue:
		return Red
	case Undefined:
		fallthrough
	default:
		panic("unexpected color")
	}
}

func (c RGB) Preview() RGB {
	switch c {
	case Red:
		return Blue
	case Green:
		return Red
	case Blue:
		return Green
	case Undefined:
		fallthrough
	default:
		panic("unexpected color")
	}
}

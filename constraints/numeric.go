package constraints

type Complex interface {
	~complex64 | ~complex128
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	SignedInt | UnsignedInt
}

type Numeric interface {
	Integer | Float | Complex
}

type OrderedNumeric interface {
	Integer | Float
}

type SignedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

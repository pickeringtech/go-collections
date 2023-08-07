package constraints

// Complex matches any complex number type (complex64 and complex128).
type Complex interface {
	~complex64 | ~complex128
}

// Float matches any float number type (float32 and float64).
type Float interface {
	~float32 | ~float64
}

// Integer matches any signed or unsigned integer type.
type Integer interface {
	SignedInt | UnsignedInt
}

// Numeric matches any non-complex numeric type (both integers and floats).
type Numeric interface {
	Integer | Float
}

// ComplexNumeric matches any numeric type (integers, floats and complex numbers).
type ComplexNumeric interface {
	Integer | Float | Complex
}

// SignedInt matches any signed integer type.
type SignedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInt matches any unsigned integer type.
type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

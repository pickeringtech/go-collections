package constraints

// Ordered matches any ordered primitive type (integers, floats and strings).
type Ordered interface {
	Integer | Float | ~string
}

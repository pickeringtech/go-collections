package constraints

type Ordered interface {
	Integer | Float | ~string
}

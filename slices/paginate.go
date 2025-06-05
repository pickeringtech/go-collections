package slices

// Paginate returns a sub-slice of the input slice based on the page index and page size.
func Paginate[T any](slice []T, pageIndex, pageSize int) []T {
	if pageIndex < 0 {
		return nil
	}

	fromIdx := pageSize * pageIndex
	toIdx := fromIdx + pageSize

	if fromIdx >= len(slice) {
		return nil
	}

	if toIdx > len(slice) {
		toIdx = len(slice)
	}

	return slice[fromIdx:toIdx]
}

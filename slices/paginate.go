package slices

// Paginate returns a sub-slice of the input slice based on the page index and page size.
//
// A non-positive pageSize (pageSize <= 0) has no meaningful page to return, so
// Paginate returns nil. Likewise, a negative pageIndex returns nil, as does a
// page that starts beyond the end of the slice - however large pageIndex and
// pageSize are, since the page bounds are computed without overflowing int.
func Paginate[T any](slice []T, pageIndex, pageSize int) []T {
	if pageIndex < 0 || pageSize <= 0 {
		return nil
	}

	// Find the page by division rather than multiplication: pageIndex*pageSize
	// can overflow int for large arguments - wrapping to a negative offset (a
	// panic) or to a wrong but valid one (the wrong page) - whereas the index of
	// the last page, (len-1)/pageSize, cannot.
	if len(slice) == 0 || pageIndex > (len(slice)-1)/pageSize {
		return nil
	}

	// Both bounds now stay within len(slice), so neither can overflow.
	fromIdx := pageIndex * pageSize
	toIdx := fromIdx + min(pageSize, len(slice)-fromIdx)

	return slice[fromIdx:toIdx]
}

package slices

// Paginate returns the page of the input slice selected by the page index and page size: the elements from
// pageIndex*pageSize up to (but excluding) the next page's start, or to the end of the slice for a partial last page.
// The page is a new slice holding copies of those elements, so it never shares the input's backing array: writing to
// or appending to it cannot affect the input.
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

	return Copy(slice[fromIdx:toIdx])
}

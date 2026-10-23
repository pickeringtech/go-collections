package slices

// Paginate returns the page of the input slice selected by the page index and page size: the elements from
// pageIndex*pageSize up to (but excluding) the next page's start, or to the end of the slice for a partial last page.
// Paginate is O(1): the page is a view of the input's elements, not a copy, with its capacity clipped to its length -
// so appending to it always allocates fresh storage and can never overwrite the input. Like any Go sub-slice, writing
// to one of its elements writes to the input's element too; use Copy for an independent slice.
//
// A non-positive pageSize (pageSize <= 0) has no meaningful page to return, so
// Paginate returns an initialised, non-nil empty slice. Likewise, a negative
// pageIndex returns an empty slice, as do an empty or nil input and a page that
// starts beyond the end of the slice - however large pageIndex and pageSize are,
// since the page bounds are computed without overflowing int.
func Paginate[T any](slice []T, pageIndex, pageSize int) []T {
	if pageIndex < 0 || pageSize <= 0 {
		return []T{}
	}

	// Find the page by division rather than multiplication: pageIndex*pageSize
	// can overflow int for large arguments - wrapping to a negative offset (a
	// panic) or to a wrong but valid one (the wrong page) - whereas the index of
	// the last page, (len-1)/pageSize, cannot.
	if len(slice) == 0 || pageIndex > (len(slice)-1)/pageSize {
		return []T{}
	}

	// Both bounds now stay within len(slice), so neither can overflow.
	fromIdx := pageIndex * pageSize
	toIdx := fromIdx + min(pageSize, len(slice)-fromIdx)

	return slice[fromIdx:toIdx:toIdx]
}

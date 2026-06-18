package slices

// Chunk splits the input into consecutive groups ("chunks") of at most size
// elements, preserving order. Every chunk holds exactly size elements except
// possibly the last, which keeps the remainder when len(input) is not an exact
// multiple of size. For example, Chunk([1 2 3 4 5], 2) is [[1 2] [3 4] [5]].
//
// If size <= 0 the request is meaningless, so the output is an initialised,
// non-nil empty slice. Likewise, if the input is empty or nil the output is an
// initialised, non-nil empty slice.
//
// Each chunk is a view into the input's backing array (no elements are copied)
// with its capacity clamped to its length, so appending to a chunk allocates a
// fresh array rather than overwriting the next chunk. The input is never
// mutated. Window is the overlapping counterpart of this fixed-step batching.
func Chunk[T any](input []T, size int) [][]T {
	output := [][]T{}
	if size <= 0 {
		return output
	}
	for i := 0; i < len(input); i += size {
		end := i + size
		if end > len(input) {
			end = len(input)
		}
		output = append(output, input[i:end:end])
	}
	return output
}

// Window returns every overlapping sub-slice ("sliding window") of width size,
// advancing one element at a time and preserving order. For example,
// Window([1 2 3 4], 2) is [[1 2] [2 3] [3 4]], yielding len(input)-size+1
// windows.
//
// If size <= 0 the request is meaningless, so the output is an initialised,
// non-nil empty slice. If size is larger than the input there is no full
// window, so the output is likewise an initialised, non-nil empty slice (this
// also covers empty or nil input). Unlike Chunk, the last window is never a
// partial one: every window has exactly size elements.
//
// Each window is a view into the input's backing array (no elements are copied)
// with its capacity clamped to its length, so appending to a window allocates a
// fresh array rather than overwriting the overlapping neighbour. The input is
// never mutated.
func Window[T any](input []T, size int) [][]T {
	output := [][]T{}
	if size <= 0 {
		return output
	}
	for i := 0; i+size <= len(input); i++ {
		end := i + size
		output = append(output, input[i:end:end])
	}
	return output
}

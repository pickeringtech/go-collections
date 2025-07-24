// Package constraints provides type constraints for building type-safe generic functions
// and APIs. These constraints extend Go's built-in comparable constraint with more
// specific type categories for numeric operations, ordering, and custom constraints.
//
// # Quick Start
//
//	import "github.com/pickeringtech/go-collections/constraints"
//
//	// Create type-safe numeric functions
//	func Sum[T constraints.Numeric](numbers []T) T {
//		var sum T
//		for _, n := range numbers {
//			sum += n
//		}
//		return sum
//	}
//
//	// Works with any numeric type
//	intSum := Sum([]int{1, 2, 3, 4, 5})           // 15
//	floatSum := Sum([]float64{1.1, 2.2, 3.3})     // 6.6
//
//	// Create ordering functions
//	func Max[T constraints.Ordered](a, b T) T {
//		if a > b { return a }
//		return b
//	}
//
//	maxInt := Max(10, 20)           // 20
//	maxString := Max("apple", "banana") // "banana"
//
// # Available Constraints
//
// Numeric Constraints:
//   - Integer: All integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr)
//   - Signed: Signed integer types (int, int8, int16, int32, int64)
//   - Unsigned: Unsigned integer types (uint, uint8, uint16, uint32, uint64, uintptr)
//   - Float: Floating point types (float32, float64)
//   - Numeric: All numeric types (Integer | Float)
//
// Ordering Constraints:
//   - Ordered: Types that support comparison operators (Numeric | string)
//
// # Common Use Cases
//
// Mathematical Functions:
//
//	func Average[T constraints.Numeric](numbers []T) T {
//		if len(numbers) == 0 {
//			return T(0)
//		}
//		sum := Sum(numbers)
//		return sum / T(len(numbers))
//	}
//
//	func Abs[T constraints.Signed](n T) T {
//		if n < 0 { return -n }
//		return n
//	}
//
// Generic Data Structures:
//
//	type MinHeap[T constraints.Ordered] struct {
//		items []T
//	}
//
//	func (h *MinHeap[T]) Push(item T) {
//		h.items = append(h.items, item)
//		h.bubbleUp(len(h.items) - 1)
//	}
//
//	func (h *MinHeap[T]) Pop() T {
//		if len(h.items) == 0 {
//			var zero T
//			return zero
//		}
//		min := h.items[0]
//		// ... heap operations
//		return min
//	}
//
// Sorting and Comparison:
//
//	func Sort[T constraints.Ordered](slice []T) {
//		sort.Slice(slice, func(i, j int) bool {
//			return slice[i] < slice[j]
//		})
//	}
//
//	func IsSorted[T constraints.Ordered](slice []T) bool {
//		for i := 1; i < len(slice); i++ {
//			if slice[i] < slice[i-1] {
//				return false
//			}
//		}
//		return true
//	}
//
// # Integration with Collections
//
// Constraints work seamlessly with the collections package:
//
//	// Numeric operations on collections
//	numbers := collections.NewList(1, 2, 3, 4, 5)
//	sum := Sum(numbers.GetAsSlice())
//
//	// Ordered operations on sets
//	scores := collections.NewSet(95, 87, 92, 78, 88)
//	sortedScores := Sort(scores.AsSlice())
//
//	// Generic functions with dicts
//	func SumValues[K comparable, V constraints.Numeric](dict collections.Dict[K, V]) V {
//		var sum V
//		dict.ForEachValue(func(value V) {
//			sum += value
//		})
//		return sum
//	}
//
// # Building Custom Constraints
//
// You can create your own constraints by combining existing ones:
//
//	// Custom constraint for types that can be zero-checked
//	type Zeroable interface {
//		comparable
//	}
//
//	func IsZero[T Zeroable](value T) bool {
//		var zero T
//		return value == zero
//	}
//
//	// Constraint for types that support arithmetic and comparison
//	type ArithmeticOrdered interface {
//		constraints.Numeric
//		constraints.Ordered
//	}
//
//	func Clamp[T ArithmeticOrdered](value, min, max T) T {
//		if value < min { return min }
//		if value > max { return max }
//		return value
//	}
//
// # Performance Considerations
//
// Generic functions with constraints have minimal runtime overhead:
//   - Type checking happens at compile time
//   - No boxing/unboxing of values
//   - Inlined for simple operations
//   - Same performance as type-specific functions
//
// Use constraints to:
//   - Build reusable, type-safe functions
//   - Create generic data structures
//   - Eliminate code duplication
//   - Catch type errors at compile time
//
// # Best Practices
//
// 1. Use the most specific constraint possible:
//    - Use Integer instead of Numeric for integer-only operations
//    - Use Signed instead of Integer for operations that need negative numbers
//
// 2. Combine constraints when needed:
//    - Create custom interfaces that combine multiple constraints
//    - Use type unions for specific type sets
//
// 3. Provide clear function signatures:
//    - Use descriptive type parameter names
//    - Document constraint requirements in function comments
//
// Start with the basic constraints (Numeric, Ordered) and create custom ones
// as your generic programming needs grow.
package constraints

# Constraints - Type-Safe Generic Programming

The `constraints` package provides type constraints for building type-safe generic functions and APIs. These constraints extend Go's built-in `comparable` with more specific type categories, enabling you to write reusable, type-safe code.

## 🚀 Quick Start

```go
import "github.com/pickeringtech/go-collections/constraints"

// Create type-safe numeric functions
func Sum[T constraints.Numeric](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}

// Works with any numeric type
intSum := Sum([]int{1, 2, 3, 4, 5})           // 15
floatSum := Sum([]float64{1.1, 2.2, 3.3})     // 6.6

// Create ordering functions
func Max[T constraints.Ordered](a, b T) T {
    if a > b { return a }
    return b
}

maxInt := Max(10, 20)                    // 20
maxString := Max("apple", "banana")      // "banana"
```

## ✨ Why Use Constraints?

**Without constraints - code duplication:**
```go
// Need separate functions for each type
func SumInts(numbers []int) int {
    var sum int
    for _, n := range numbers { sum += n }
    return sum
}

func SumFloats(numbers []float64) float64 {
    var sum float64
    for _, n := range numbers { sum += n }
    return sum
}

// ... more functions for each numeric type
```

**With constraints - one generic function:**
```go
// Single function works for all numeric types
func Sum[T constraints.Numeric](numbers []T) T {
    var sum T
    for _, n := range numbers { sum += n }
    return sum
}

// Type-safe and reusable
intSum := Sum([]int{1, 2, 3})
floatSum := Sum([]float64{1.1, 2.2, 3.3})
```

## 📋 Available Constraints

### 🔢 Numeric Constraints

#### Integer - All Integer Types
```go
// Includes: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr
func CountBits[T constraints.Integer](n T) int {
    count := 0
    for n > 0 {
        count++
        n >>= 1
    }
    return count
}

bits := CountBits(42)        // Works with int
bits = CountBits(uint8(42))  // Works with uint8
bits = CountBits(int64(42))  // Works with int64
```

#### Signed - Signed Integer Types
```go
// Includes: int, int8, int16, int32, int64
func Abs[T constraints.Signed](n T) T {
    if n < 0 { return -n }
    return n
}

result := Abs(-42)      // 42 (works with signed integers)
result = Abs(int8(-5))  // 5 (works with int8)
// result = Abs(uint(5)) // Compile error - uint is not signed
```

#### Unsigned - Unsigned Integer Types
```go
// Includes: uint, uint8, uint16, uint32, uint64, uintptr
func IsPowerOfTwo[T constraints.Unsigned](n T) bool {
    return n > 0 && (n&(n-1)) == 0
}

isPower := IsPowerOfTwo(uint(8))   // true
isPower = IsPowerOfTwo(uint16(7))  // false
// isPower = IsPowerOfTwo(-8)      // Compile error - negative not allowed
```

#### Float - Floating Point Types
```go
// Includes: float32, float64
func Round[T constraints.Float](n T, decimals int) T {
    multiplier := T(math.Pow(10, float64(decimals)))
    return T(math.Round(float64(n)*float64(multiplier))) / multiplier
}

rounded := Round(3.14159, 2)        // 3.14 (float64)
rounded32 := Round(float32(3.14159), 2) // 3.14 (float32)
```

#### Numeric - All Numeric Types
```go
// Includes: Integer | Float (all numeric types)
func Average[T constraints.Numeric](numbers []T) T {
    if len(numbers) == 0 {
        return T(0)
    }
    sum := Sum(numbers)
    return sum / T(len(numbers))
}

avgInt := Average([]int{1, 2, 3, 4, 5})           // 3
avgFloat := Average([]float64{1.5, 2.5, 3.5})     // 2.5
```

### 📊 Ordering Constraints

#### Ordered - Types Supporting Comparison
```go
// Includes: Numeric | string (all comparable types with <, >, etc.)
func Sort[T constraints.Ordered](slice []T) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i] < slice[j]
    })
}

numbers := []int{3, 1, 4, 1, 5}
Sort(numbers)  // [1, 1, 3, 4, 5]

words := []string{"banana", "apple", "cherry"}
Sort(words)    // ["apple", "banana", "cherry"]

func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}

minNum := Min(10, 20)           // 10
minStr := Min("zebra", "apple") // "apple"
```

## 🛠️ Building Generic Data Structures

### Generic Stack
```go
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// Type-safe usage
intStack := NewStack[int]()
intStack.Push(42)
value, ok := intStack.Pop() // value is int

stringStack := NewStack[string]()
stringStack.Push("hello")
text, ok := stringStack.Pop() // text is string
```

### Generic Priority Queue
```go
type PriorityQueue[T constraints.Ordered] struct {
    items []T
}

func NewPriorityQueue[T constraints.Ordered]() *PriorityQueue[T] {
    return &PriorityQueue[T]{items: make([]T, 0)}
}

func (pq *PriorityQueue[T]) Push(item T) {
    pq.items = append(pq.items, item)
    // Bubble up logic here
}

func (pq *PriorityQueue[T]) Pop() (T, bool) {
    if len(pq.items) == 0 {
        var zero T
        return zero, false
    }

    min := pq.items[0]
    // Heap operations here
    return min, true
}

// Works with any ordered type
intPQ := NewPriorityQueue[int]()
intPQ.Push(5)
intPQ.Push(2)
intPQ.Push(8)
min, _ := intPQ.Pop() // 2

stringPQ := NewPriorityQueue[string]()
stringPQ.Push("zebra")
stringPQ.Push("apple")
first, _ := stringPQ.Pop() // "apple"
```

### Generic Binary Search Tree
```go
type BST[T constraints.Ordered] struct {
    value T
    left  *BST[T]
    right *BST[T]
}

func NewBST[T constraints.Ordered](value T) *BST[T] {
    return &BST[T]{value: value}
}

func (bst *BST[T]) Insert(value T) {
    if value < bst.value {
        if bst.left == nil {
            bst.left = NewBST(value)
        } else {
            bst.left.Insert(value)
        }
    } else {
        if bst.right == nil {
            bst.right = NewBST(value)
        } else {
            bst.right.Insert(value)
        }
    }
}

func (bst *BST[T]) Contains(value T) bool {
    if value == bst.value {
        return true
    } else if value < bst.value && bst.left != nil {
        return bst.left.Contains(value)
    } else if value > bst.value && bst.right != nil {
        return bst.right.Contains(value)
    }
    return false
}

// Type-safe usage
intBST := NewBST(10)
intBST.Insert(5)
intBST.Insert(15)
found := intBST.Contains(5) // true

stringBST := NewBST("middle")
stringBST.Insert("apple")
stringBST.Insert("zebra")
found = stringBST.Contains("apple") // true
```

## 🌟 Real-World Examples

### Mathematical Operations
```go
// Generic mathematical functions
func Factorial[T constraints.Integer](n T) T {
    if n <= 1 {
        return 1
    }
    return n * Factorial(n-1)
}

func GCD[T constraints.Integer](a, b T) T {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func Clamp[T constraints.Ordered](value, min, max T) T {
    if value < min { return min }
    if value > max { return max }
    return value
}

// Usage
fact := Factorial(5)           // 120
gcd := GCD(48, 18)            // 6
clamped := Clamp(15, 10, 20)  // 15
clampedStr := Clamp("m", "a", "z") // "m"
```

### Data Processing
```go
// Generic data processing functions
func FindMax[T constraints.Ordered](slice []T) (T, bool) {
    if len(slice) == 0 {
        var zero T
        return zero, false
    }

    max := slice[0]
    for _, item := range slice[1:] {
        if item > max {
            max = item
        }
    }
    return max, true
}

func Unique[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    result := make([]T, 0)

    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}

func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
    var trueSlice, falseSlice []T

    for _, item := range slice {
        if predicate(item) {
            trueSlice = append(trueSlice, item)
        } else {
            falseSlice = append(falseSlice, item)
        }
    }
    return trueSlice, falseSlice
}

// Usage
numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}
max, found := FindMax(numbers)  // 9, true

unique := Unique(numbers)       // [3, 1, 4, 5, 9, 2, 6]

evens, odds := Partition(numbers, func(n int) bool { return n%2 == 0 })
// evens: [4, 2, 6], odds: [3, 1, 1, 5, 9]
```

### Configuration and Validation
```go
// Generic configuration with validation
type Config[T constraints.Numeric] struct {
    MinValue T
    MaxValue T
    Default  T
}

func (c Config[T]) Validate(value T) T {
    if value < c.MinValue {
        return c.MinValue
    }
    if value > c.MaxValue {
        return c.MaxValue
    }
    return value
}

func (c Config[T]) ValidateWithDefault(value *T) T {
    if value == nil {
        return c.Default
    }
    return c.Validate(*value)
}

// Usage
intConfig := Config[int]{MinValue: 0, MaxValue: 100, Default: 50}
validated := intConfig.Validate(150)  // 100 (clamped to max)

floatConfig := Config[float64]{MinValue: 0.0, MaxValue: 1.0, Default: 0.5}
validated = floatConfig.Validate(1.5) // 1.0 (clamped to max)
```

## 🔗 Integration with Collections

Constraints work seamlessly with the collections package:

```go
// Numeric operations on collections
func SumCollection[T constraints.Numeric](list collections.List[T]) T {
    var sum T
    list.ForEach(func(value T) {
        sum += value
    })
    return sum
}

func MaxInSet[T constraints.Ordered](set collections.Set[T]) (T, bool) {
    slice := set.AsSlice()
    return FindMax(slice)
}

func SumDictValues[K comparable, V constraints.Numeric](dict collections.Dict[K, V]) V {
    var sum V
    dict.ForEachValue(func(value V) {
        sum += value
    })
    return sum
}

// Usage with collections
numbers := collections.NewList(1, 2, 3, 4, 5)
total := SumCollection(numbers) // 15

scores := collections.NewSet(95, 87, 92, 78, 88)
highest, found := MaxInSet(scores) // 95, true

inventory := collections.NewDict(
    collections.Pair[string, int]{Key: "apples", Value: 50},
    collections.Pair[string, int]{Key: "oranges", Value: 30},
)
totalItems := SumDictValues(inventory) // 80
```

## 🎯 Building Custom Constraints

### Combining Existing Constraints
```go
// Custom constraint for arithmetic operations
type ArithmeticOrdered interface {
    constraints.Numeric
    constraints.Ordered
}

func Lerp[T ArithmeticOrdered](a, b T, t float64) T {
    return a + T(float64(b-a)*t)
}

// Custom constraint for types that can be zero-checked
type Zeroable interface {
    comparable
}

func IsZero[T Zeroable](value T) bool {
    var zero T
    return value == zero
}

func IsNonZero[T Zeroable](value T) bool {
    return !IsZero(value)
}

// Usage
interpolated := Lerp(10, 20, 0.5)    // 15
isZero := IsZero(0)                   // true
isZero = IsZero("")                   // true
isZero = IsZero(42)                   // false
```

### Creating Domain-Specific Constraints
```go
// Constraint for ID types
type ID interface {
    constraints.Integer
    ~int | ~int64 | ~uint | ~uint64
}

func ValidateID[T ID](id T) bool {
    return id > 0
}

// Constraint for percentage types
type Percentage interface {
    constraints.Float
    ~float32 | ~float64
}

func ValidatePercentage[T Percentage](p T) bool {
    return p >= 0 && p <= 100
}

// Custom types that satisfy constraints
type UserID int64
type Score float64

// Usage
userID := UserID(12345)
valid := ValidateID(userID) // true

score := Score(85.5)
valid = ValidatePercentage(score) // true
```

## 📊 Performance Guide

### Compile-Time vs Runtime

**Constraints are compile-time only:**
```go
// This generates the same assembly as non-generic version
func AddInts(a, b int) int { return a + b }
func Add[T constraints.Numeric](a, b T) T { return a + b }

// Both compile to identical machine code
result1 := AddInts(5, 10)
result2 := Add(5, 10)
```

**No runtime overhead:**
- Type checking happens at compile time
- No boxing/unboxing of values
- Functions are inlined for simple operations
- Same performance as type-specific functions

### Best Practices for Performance

```go
// ✅ Good: Use specific constraints
func ProcessIntegers[T constraints.Integer](data []T) { ... }

// ❌ Avoid: Overly broad constraints when specific ones work
func ProcessIntegers[T constraints.Numeric](data []T) { ... }

// ✅ Good: Inline simple operations
func Add[T constraints.Numeric](a, b T) T {
    return a + b  // Will be inlined
}

// ✅ Good: Use constraints for type safety, not performance
func SafeDivide[T constraints.Float](a, b T) T {
    if b == 0 {
        return 0
    }
    return a / b
}
```

## 🎯 Best Practices

### 1. 🎨 Use the Most Specific Constraint
```go
// ✅ Good: Use specific constraint
func Abs[T constraints.Signed](n T) T {
    if n < 0 { return -n }
    return n
}

// ❌ Avoid: Overly broad constraint
func Abs[T constraints.Numeric](n T) T {
    // Won't work with unsigned types anyway
}
```

### 2. 🔧 Combine Constraints When Needed
```go
// ✅ Good: Create meaningful combinations
type Comparable[T any] interface {
    comparable
}

type OrderedComparable[T any] interface {
    constraints.Ordered
    comparable
}

func UniqueAndSort[T OrderedComparable[T]](slice []T) []T {
    unique := Unique(slice)
    Sort(unique)
    return unique
}
```

### 3. 📝 Document Constraint Requirements
```go
// ✅ Good: Clear documentation
// ProcessNumericData processes a slice of numeric values.
// T must be a numeric type (integer or float).
func ProcessNumericData[T constraints.Numeric](data []T) T {
    // Implementation
}

// ✅ Good: Explain why constraint is needed
// BinarySearch requires ordered types to perform comparisons.
// T must support <, >, and == operators.
func BinarySearch[T constraints.Ordered](slice []T, target T) int {
    // Implementation
}
```

### 4. 🧪 Test with Multiple Types
```go
func TestSum(t *testing.T) {
    // Test with different numeric types
    intResult := Sum([]int{1, 2, 3})
    assert.Equal(t, 6, intResult)

    floatResult := Sum([]float64{1.1, 2.2, 3.3})
    assert.InDelta(t, 6.6, floatResult, 0.001)

    int64Result := Sum([]int64{1, 2, 3})
    assert.Equal(t, int64(6), int64Result)
}
```

## 🚀 Quick Reference

### Available Constraints
```go
constraints.Integer    // All integer types
constraints.Signed     // Signed integers only
constraints.Unsigned   // Unsigned integers only
constraints.Float      // Floating point types
constraints.Numeric    // All numeric types
constraints.Ordered    // Types supporting <, >, ==
```

### Common Patterns
```go
// Mathematical functions
func Sum[T constraints.Numeric]([]T) T
func Max[T constraints.Ordered](T, T) T
func Abs[T constraints.Signed](T) T

// Data structures
type Stack[T any] struct { ... }
type PriorityQueue[T constraints.Ordered] struct { ... }

// Custom constraints
type MyConstraint interface {
    constraints.Numeric
    ~int | ~float64  // Specific types only
}
```

Start with the basic constraints (`Numeric`, `Ordered`) and create custom ones as your generic programming needs grow. Remember: constraints are about type safety, not performance!

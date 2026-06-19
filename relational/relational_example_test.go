package relational_test

import (
	"fmt"
	"sort"

	"github.com/pickeringtech/go-collections/relational"
	"github.com/pickeringtech/go-collections/stats"
)

type saleRow struct {
	Dept   string
	Amount float64
}

// Example_quickStart mirrors the doc.go Quick Start: GROUP BY + aggregate, then
// a join, so the package documentation is guaranteed to compile and run.
func Example_quickStart() {
	orders := []saleRow{
		{"books", 10}, {"books", 30}, {"toys", 5},
	}

	byDept := relational.GroupBy(orders, func(o saleRow) string { return o.Dept })
	avg := relational.AggregateBy(byDept,
		func(o saleRow) float64 { return o.Amount },
		stats.Mean[float64],
	)
	fmt.Printf("books avg=%.0f toys avg=%.0f\n", avg["books"], avg["toys"])

	type label struct {
		Dept string
		Name string
	}
	labels := []label{{"books", "Books Dept"}, {"toys", "Toys Dept"}}
	joined := relational.InnerJoin(orders, labels,
		func(o saleRow) string { return o.Dept },
		func(l label) string { return l.Dept },
	)
	fmt.Printf("joined rows=%d\n", len(joined))
	// Output:
	// books avg=20 toys avg=5
	// joined rows=3
}

func ExampleGroupBy() {
	nums := []int{1, 2, 3, 4, 5}
	groups := relational.GroupBy(nums, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	fmt.Println("odd:", groups["odd"])
	fmt.Println("even:", groups["even"])
	// Output:
	// odd: [1 3 5]
	// even: [2 4]
}

func ExampleGroupBySeq() {
	// A pull sequence (here over a slice, but it could be a DB cursor) is
	// grouped without first collecting it into a slice.
	seq := func(yield func(int) bool) {
		for _, n := range []int{1, 2, 3, 4} {
			if !yield(n) {
				return
			}
		}
	}
	groups := relational.GroupBySeq(seq, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	fmt.Println("odd:", groups["odd"])
	fmt.Println("even:", groups["even"])
	// Output:
	// odd: [1 3]
	// even: [2 4]
}

func ExampleAggregateBy() {
	type sale struct {
		Region string
		Amount float64
	}
	sales := []sale{
		{"EU", 100}, {"EU", 200}, {"US", 50},
	}
	byRegion := relational.GroupBy(sales, func(s sale) string { return s.Region })
	totals := relational.AggregateBy(byRegion,
		func(s sale) float64 { return s.Amount },
		stats.Sum[float64],
	)
	fmt.Printf("EU=%.0f US=%.0f\n", totals["EU"], totals["US"])
	// Output:
	// EU=300 US=50
}

func ExampleRightJoin() {
	type user struct{ ID int }
	type post struct{ UserID int }
	users := []user{{1}}
	posts := []post{{1}, {2}}

	pairs := relational.RightJoin(users, posts,
		func(u user) int { return u.ID },
		func(p post) int { return p.UserID },
	)
	for _, p := range pairs {
		fmt.Printf("post.user=%d hasUser=%v\n", p.Right.UserID, p.LeftOK)
	}
	// Output:
	// post.user=1 hasUser=true
	// post.user=2 hasUser=false
}

func ExampleFullOuterJoin() {
	type a struct{ K int }
	type b struct{ K int }
	as := []a{{1}, {2}}
	bs := []b{{1}, {3}}

	pairs := relational.FullOuterJoin(as, bs,
		func(x a) int { return x.K },
		func(y b) int { return y.K },
	)
	for _, p := range pairs {
		fmt.Printf("L(%d,%v) R(%d,%v)\n", p.Left.K, p.LeftOK, p.Right.K, p.RightOK)
	}
	// Output:
	// L(1,true) R(1,true)
	// L(2,true) R(0,false)
	// L(0,false) R(3,true)
}

func ExampleCountBy() {
	words := []string{"apple", "avocado", "banana", "cherry"}
	counts := relational.CountBy(words, func(w string) byte { return w[0] })
	fmt.Println("a:", counts['a'])
	fmt.Println("b:", counts['b'])
	fmt.Println("c:", counts['c'])
	// Output:
	// a: 2
	// b: 1
	// c: 1
}

func ExampleAggregate() {
	groups := map[string][]int{
		"a": {1, 2, 3},
		"b": {10, 20},
	}
	totals := relational.Aggregate(groups, stats.Sum[int])
	fmt.Println("a:", totals["a"])
	fmt.Println("b:", totals["b"])
	// Output:
	// a: 6
	// b: 30
}

func ExampleInnerJoin() {
	type user struct {
		ID   int
		Name string
	}
	type post struct {
		UserID int
		Title  string
	}
	users := []user{{1, "alice"}, {2, "bob"}}
	posts := []post{{1, "hello"}, {1, "world"}}

	pairs := relational.InnerJoin(users, posts,
		func(u user) int { return u.ID },
		func(p post) int { return p.UserID },
	)
	for _, p := range pairs {
		fmt.Printf("%s -> %s\n", p.Left.Name, p.Right.Title)
	}
	// Output:
	// alice -> hello
	// alice -> world
}

func ExampleLeftJoin() {
	type user struct{ ID int }
	type post struct{ UserID int }
	users := []user{{1}, {2}}
	posts := []post{{1}}

	pairs := relational.LeftJoin(users, posts,
		func(u user) int { return u.ID },
		func(p post) int { return p.UserID },
	)
	for _, p := range pairs {
		fmt.Printf("user=%d hasPost=%v\n", p.Left.ID, p.RightOK)
	}
	// Output:
	// user=1 hasPost=true
	// user=2 hasPost=false
}

func ExamplePivot() {
	type cell struct {
		Quarter string
		Region  string
		Sales   int
	}
	rows := []cell{
		{"Q1", "EU", 10}, {"Q1", "US", 20}, {"Q2", "EU", 30},
	}
	wide := relational.Pivot(rows,
		func(c cell) string { return c.Quarter },
		func(c cell) string { return c.Region },
		func(c cell) int { return c.Sales },
	)
	fmt.Println("Q1 EU:", wide["Q1"]["EU"])
	fmt.Println("Q1 US:", wide["Q1"]["US"])
	fmt.Println("Q2 EU:", wide["Q2"]["EU"])
	// Output:
	// Q1 EU: 10
	// Q1 US: 20
	// Q2 EU: 30
}

func ExampleUnpivot() {
	wide := map[string]map[string]int{
		"Q1": {"EU": 10, "US": 20},
	}
	cells := relational.Unpivot(wide)
	// Sort for a stable example output (map order is randomised).
	sort.Slice(cells, func(i, j int) bool { return cells[i].Col < cells[j].Col })
	for _, c := range cells {
		fmt.Printf("%s/%s=%d\n", c.Row, c.Col, c.Value)
	}
	// Output:
	// Q1/EU=10
	// Q1/US=20
}

func ExamplePartition() {
	nums := []int{1, 2, 3, 4, 5, 6}
	evens, odds := relational.Partition(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)
	fmt.Println("odds:", odds)
	// Output:
	// evens: [2 4 6]
	// odds: [1 3 5]
}

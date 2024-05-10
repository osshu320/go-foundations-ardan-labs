package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"
)

func findExample() {
	target := 4
	a := []int{1, 2, 3, 4, 5, 6}
	i, f := sort.Find(len(a), func(x int) int {
		if target > a[x] {
			return 1
		} else {
			return 0
		}
	})

	if f {
		fmt.Printf("found %d at %d", target, i)
	} else {
		fmt.Println("Not found")
		fmt.Println(a)
	}
}

func searchExample() {
	// x := 6
	// a := []int{2, 4, 6, 8, 10}
	x := 5.0
	a := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}

	ind := sort.Search(len(a), func(i int) bool {
		return a[i] >= x
	})

	if ind < len(a) && a[ind] == x {
		fmt.Println("x is present at", ind)
	} else if ind < len(a) {
		fmt.Println("just greater that x is present at ind", ind)
	} else {
		fmt.Println("nothing greater than x present")
	}
}

func sortExamples() {
	names := []string{"Bob", "Alice", "VERA"}
	slices.SortFunc(names, func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})
	fmt.Println(names)

	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Gopher", 13},
		{"Alice", 55},
		{"Bob", 24},
		{"Alice", 20},
	}

	slices.SortFunc(people, func(a, b Person) int {
		if n := cmp.Compare(a.Name, b.Name); n != 0 {
			return n
		}
		return cmp.Compare(a.Age, b.Age)
	})

	fmt.Println(people)

	s := []int{1, 2, 3, 4, 5}
	t := sort.Reverse(sort.IntSlice(s))
	fmt.Println(s)
	fmt.Println(t)
}

func main() {

}

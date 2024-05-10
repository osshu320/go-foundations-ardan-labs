package main

import "fmt"

func SubOneFromLen(slice []int) []int {
	slice = slice[0 : len(slice)-1]
	return slice
}

func slices_go_main() {
	a := []int{1, 2, 3, 4, 5}
	fmt.Println("a", len(a), a)
	b := SubOneFromLen(a)
	fmt.Println("a", len(a), a)
	fmt.Println("b", len(b), b)
}

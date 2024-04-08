package main

import (
	"fmt"
	"unicode/utf8"
)

func banner(text string, width int) {
	padding := (width - utf8.RuneCountInString(text)) / 2
	// padding := (width - len(text)) / 2 // BUG : len is in bytes
	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}
	fmt.Println(text)
	for i := 0; i < width; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func isPallindrome(s string) bool {
	n := len(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func doSomeFun() {
	s := "G☺"

	fmt.Println("len(s): ", len(s))
	fmt.Println("utf8.RuneCountInString(s):", utf8.RuneCountInString(s))

	fmt.Println("s[0]: ", s[0])
	fmt.Println("s[1]: ", s[1])

	fmt.Println("[]rune(s)[0]: ", []rune(s)[0])
	fmt.Println("[]rune(s)[1]: ", []rune(s)[1])

	fmt.Println("Iterating over s:")
	for i, x := range s {
		fmt.Println(i, x)
	}
}

// len(s):  4
// utf8.RuneCountInString(s): 2
// s[0]:  71
// s[1]:  226
// []rune(s)[0]:  71
// []rune(s)[1]:  9786
// Iterating over s:
// 0 71
// 1 9786

func main() {
	// banner("GO", 6)
	// str := "G☺"
	// banner(str, 6)
	// fmt.Println(len(str))

	// for i, r := range str {
	// 	fmt.Println(i, r)
	// 	fmt.Printf("%c of type %T\n", r, r)
	// }

	// b := str[0]
	// c := str[1]
	// d := str[2]
	// fmt.Printf("%c of type %T\n", b, b)
	// fmt.Printf("%c of type %T\n", c, c)
	// fmt.Printf("%c of type %T\n", d, d)

	// p := []rune(str)[0]
	// q := []rune(str)[1]
	// fmt.Printf("%c of type %T\n", p, p)
	// fmt.Printf("%c of type %T\n", q, q)

	// s := "G☺"
	// fmt.Println(len(s))                    // Prints: 4
	// fmt.Println(utf8.RuneCountInString(s)) // Prints: 2

	doSomeFun()

	// fmt.Println(isPallindrome("g"))
	// fmt.Println(isPallindrome("go"))
	// fmt.Println(isPallindrome("gog"))
	// fmt.Println(isPallindrome("gogo"))
}

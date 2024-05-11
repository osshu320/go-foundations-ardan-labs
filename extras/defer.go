package main

import (
	"fmt"
	"log"
)

func a() {
	i := 0
	defer fmt.Println(i)
	i++
	return
}

func b() {
	for i := 0; i < 4; i++ {
		defer fmt.Println(i)
	}
}

func c() (i int) {
	defer func() { i++ }()
	log.Println("function body i", i)
	return 1
}

func g(i int) {
	if i > 3 {
		log.Println("Panicking")
		panic(fmt.Sprintf("%v", i))
	}

	defer log.Println("Defer in g", i)
	log.Println("Printing in g", i)
	g(i + 1)
}

func f() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	log.Println("Calling g.")
	g(0)
	log.Println("Returned Normally from g.")
}

/*
printing in g 0
printing in g 1
printing in g 2
printing in g 3
panicking
printing in g 0
defer in g 3
defer in g 2
defer in g 1
defer in g 0
recovered in f 4
returned normally from f
*/

func h() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	panic("panicking intentiionally man!")
}

func defer_go_main() {
	// a()
	// b()
	// log.Println(c())

	// f()
	// log.Println("Returned normall from f.")

	h()
	log.Println("Normal return from h")
}

package main

import (
	"log"
	"sync"
)

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func pipelines_main() {
	// c := gen(2, 3, 4)
	// out := sq(c)

	// for x := range out {
	// 	log.Println(x)
	// }

	// out := sq(gen(2, 3, 4))
	// for x := range out {
	// 	log.Println(x)
	// }

	// for x := range sq(gen(2, 3, 4)) {
	// 	log.Println(x)
	// }

	// for x := range sq(sq(gen(2, 3, 4))) {
	// 	log.Println(x)
	// }

	// in := gen(2, 3, 4, 5, 6, 7)
	// c1 := sq(in)
	// c2 := sq(in)

	// for n := range merge(c1, c2) {
	// 	log.Println(n)
	// }

	done := make(chan struct{})
	defer close(done)

	in := gen(done, 2, 3)
	c1 := sq(done, in)
	c2 := sq(done, in)

	out := merge(done, c1, c2)
	log.Println(<-out)

}

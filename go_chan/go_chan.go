package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	go fmt.Println("goroutine")
	fmt.Println("main")

	for i := 0; i < 3; i++ {
		/* BUG: All goroutines use same i from the for loop
		go func() {
			fmt.Println(i)
		}()
		*/

		/* Fix 1: Use a parameter */
		// go func(i int) {
		// 	fmt.Println(i)
		// }(i)

		/* Fix 2: Use a loop body variable */
		i := i
		go func() {
			fmt.Println(i)
		}()
	}

	time.Sleep(10 * time.Millisecond)

	ch := make(chan string)
	go func() {
		ch <- "hi" // send
	}()
	msg := <-ch // receive
	fmt.Println(msg)

	go func() {
		for i := 0; i < 3; i++ {
			msg := fmt.Sprintf("message #%d", i+1)
			ch <- msg
		}
		close(ch)
	}()

	for msg := range ch {
		fmt.Println("got:", msg)
	}

	msg, ok := <-ch // ch is closed
	fmt.Printf("closed: %#v (ok=%v)\n", msg, ok)

	// ch <- "hi" // ch is closed -> panic

	vals := []int{3, 2, 4, 1, 5}
	sortedVals := sleepSort(vals)
	log.Println(sortedVals)
}

/*
For every value "n" in values, spin a goroutine that will
- sleep "n" milliseconds
- Send "n" over a channel

In the function body, collect values from the channel to a slice and return it
*/

// Sleepsort is an example of FAN-OUT
func sleepSort(values []int) []int {
	ch := make(chan int)
	for _, x := range values {
		go func(ch chan<- int, xx int) {
			time.Sleep(time.Duration(xx) * time.Millisecond)
			ch <- xx
		}(ch, x)
	}
	// ans := make([]int, len(values))
	// for i := 0; i < len(values); i++ {
	// 	ans[i] = <-ch
	// }
	var ans []int
	for range values {
		ans = append(ans, <-ch)
	}
	return ans
}

/* Channel Semantics
- send & receive will block until opposite operation (*)
- receive from closed channel will return zero value without blocking
- send to a closed channel will panic
- closing a closed channel will panic
- send/receive to a nil channel will block forever
*/

package main

import (
	"fmt"
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
}

package main

import (
	"context"
	"log"
	"time"
)

func main() {
	ch1, ch2 := make(chan int), make(chan int)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1 <- 1
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		ch2 <- 2
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	select {
	case val := <-ch1:
		log.Println("ch1:", val)
	case val := <-ch2:
		log.Println("ch2:", val)
	// case <-time.After(5 * time.Millisecond):
	case <-ctx.Done():
		log.Println("timeout")
	}

	// select {} // blocks forever without consuming CPU
	// for {} // this consumes a lot of CPU
}

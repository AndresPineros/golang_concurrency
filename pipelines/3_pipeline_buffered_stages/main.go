package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Buffered channels is the last optimization that should be considered
because they can hide deadlocks and other concurrency issues.

Little's Law shows that buffering doesn't necessarily improve throughput,
it just reduces the amount of time the system is in a blocking state.
*/

func main() {
	generator := func(done chan interface{}) <-chan int {
		intStream := make(chan int, 30) // Buffered to make stages independent
		count := 0
		go func() {
			defer close(intStream)
			for {
				select {
				case <-time.After(200 * time.Millisecond):
					count++
					intStream <- int(rand.Uint32())
				case <-done:
					fmt.Println("Cancelled generator stage. Produced", count, "integers")
					return
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int, 30) // Buffered to make stages independent
		count := 0
		go func() {
			defer close(multipliedStream)
			for i := range intStream { // range makes sure we process all the elements in the channel.
				time.Sleep(300 * time.Millisecond)
				count++
				multipliedStream <- i * multiplier
			}
			fmt.Println("Source stream closed for multiply stage. Processed", count, "numbers and exited.")
		}()
		return multipliedStream
	}
	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int, 30) // Buffered to make stages independent
		count := 0
		go func() {
			defer close(addedStream)
			for i := range intStream {
				time.Sleep(500 * time.Millisecond)
				count++
				addedStream <- i + additive
			}
			fmt.Println("Source stream closed for add stage. Processed", count, "numbers and exited.")
		}()
		return addedStream
	}
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Shutting down generator...")
		done <- 1
	}()
	intStream := generator(done)
	pipeline := multiply(nil, add(nil, multiply(nil, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}

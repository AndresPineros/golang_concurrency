package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
A stage can only start working on the next item
after the next stage receives the value.

When the intStream is closed, all the remaining elements in the pipeline are processed.
*/

func main() {
	generator := func(done chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for {
				select {
				case <-time.After(200 * time.Millisecond):
					intStream <- int(rand.Uint32())
				case <-done:
					return
				}
			}
		}()
		return intStream
	}

	multiply := func(intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				multipliedStream <- i * multiplier
			}
		}()
		return multipliedStream
	}

	add := func(intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				addedStream <- i + additive
			}
		}()
		return addedStream
	}
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done)
	pipeline := multiply(add(multiply(intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}

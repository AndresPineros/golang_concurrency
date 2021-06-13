package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
A stage can only start working on the next item
after the next stage receives the value.
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
			fmt.Println("Multiply worker finished.")
		}()
		return multipliedStream
	}

	add := func(intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		fanout_workers := 4
		wg := sync.WaitGroup{}
		wg.Add(fanout_workers)
		go func() { // Goroutine that waits until all workers are done and closes channel
			wg.Wait()
			close(addedStream)
		}()
		for i := 0; i < fanout_workers; i++ {
			go func() {
				for i := range intStream { // we have 4 workers reading from this intStream. This is the fanout
					addedStream <- i + additive // here we fan-in back to the next stream
				}
				fmt.Println("Add worker finished.")
				wg.Done()
			}()
		}
		return addedStream
	}
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		done <- 1
	}()
	intStream := generator(done)
	pipeline := multiply(add(multiply(intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}

	time.Sleep(5 * time.Second) // just to wait for all fmt.Println
}

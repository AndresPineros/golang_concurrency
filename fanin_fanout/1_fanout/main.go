package main

import (
	"fmt"
	"time"
)

func fanout(workers int, // The amount of workers
	inputStream chan interface{}, // The input stream that will be fanned out.
	fun func(interface{}) interface{}, // The function that will be executed by each worker
) []chan interface{} {
	channels := make([]chan interface{}, workers)
	for i := 0; i < workers; i++ {
		c := make(chan interface{})
		channels[i] = c
		go func(chan interface{}) {
			defer close(c) // avoid leaking channels
			for elem := range inputStream {
				c <- fun(elem)
			}
		}(c)
	}
	return channels
}

func main() {

	workers := 4
	inputStream := make(chan interface{})
	go func() {
		defer close(inputStream) // avoid leaking channels
		for i := 0; i < 100; i++ {
			inputStream <- i
		}
	}()
	adder := func(number int, adder int) int {
		return number + adder
	}

	add5fun := func(i interface{}) interface{} { // wrapper to match the fanout() parameter function
		return adder(5, i.(int))
	}

	channels := fanout(workers, inputStream, add5fun)

	for i := 0; i < workers; i++ {
		go func(pos int) {
			for elem := range channels[pos] {
				fmt.Println("Got element", elem, "from channel", pos)
			}
			fmt.Println("Closed channel read", pos)
		}(i)
	}

	time.Sleep(5 * time.Second) // too lazy to use WaitGroup
}

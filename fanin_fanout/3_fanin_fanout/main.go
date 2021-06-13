package main

import (
	"fmt"
	"sync"
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

func fanin(channels []chan interface{}) chan interface{} {
	outputStream := make(chan interface{})
	wg := sync.WaitGroup{}
	wg.Add(len(channels))
	for i := 0; i < len(channels); i++ {
		c := channels[i]
		go func(c chan interface{}) {
			for elem := range c {
				outputStream <- elem
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(outputStream)
	}()
	return outputStream
}

func main() {
	/*
		Joining the functions for generic fan-out and fan-in...
	*/

	inputStream := make(chan interface{})
	go func() {
		defer close(inputStream) // avoid leaking channels
		for i := 0; i < 10; i++ {
			inputStream <- i
		}
	}()

	workers := 4
	fannedOut := fanout(workers, inputStream, func(i interface{}) interface{} {
		return 10 + i.(int) // This is the task that will be fannedout and performed on the inputStream
	})
	fannedIn := fanin(fannedOut)

	for elem := range fannedIn {
		fmt.Println(elem)
	}

}

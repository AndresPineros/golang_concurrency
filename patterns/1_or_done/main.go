package main

import (
	"fmt"
	"time"
)

func orDone(done, c <-chan interface{}) <-chan interface{} {
	intermediateChannel := make(chan interface{})
	go func() {
		defer close(intermediateChannel)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case intermediateChannel <- v:
				case <-done:
				}
			}
		}
	}()
	return intermediateChannel
}

func main() {
	done := make(chan interface{})
	inputStream := make(chan interface{})
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(1 * time.Second)
			inputStream <- i
		}
	}()
	go func() {
		time.Sleep(5 * time.Second)
		done <- 1 // will trigger done after 5 seconds
	}()
	/*
		The or-done function allows hiding the complexity of finishing consuming a channel when a done
		channel is passed.

		BUT it adds latency because we're using an itermediate channel, and channels are slow.
	*/
	for elem := range orDone(done, inputStream) {
		fmt.Println(elem)
	}
}

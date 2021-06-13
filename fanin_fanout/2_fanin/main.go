package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	fanin := func(channels []chan interface{}) chan interface{} {
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

	workers := 4
	channels := make([]chan interface{}, workers)
	for i := 0; i < workers; i++ {
		c := make(chan interface{})
		channels[i] = c
		go func(c chan interface{}, chanid int) {
			defer close(c)
			for i := 0; i < 5; i++ {
				c <- fmt.Sprintf("Channel %v. Value %v", chanid, i)
			}
		}(c, i)
	}
	stream := fanin(channels) // Everything is fanned into this channel
	for e := range stream {
		fmt.Println(e)
	}

	time.Sleep(2 * time.Second)
}

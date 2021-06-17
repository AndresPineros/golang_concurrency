package main

import (
	"context"
	"fmt"
	"sync"
)

/*
context.Context and Done channels both help with cancelling goroutines BUT
context also has timeout functionalities.
*/

func context_broadcast() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		<-ctx.Done()
		fmt.Println("First goroutine exited due to context")
		wg.Done()
	}()
	go func() {
		<-ctx.Done()
		fmt.Println("Second goroutine exited due to context")
		wg.Done()
	}()
	cancel()
	wg.Wait()
}

func done_channel() {
	done := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		<-done
		fmt.Println("First goroutine exited due to done channel")
		wg.Done()
	}()
	go func() {
		<-done
		fmt.Println("Second goroutine exited due to done channel")
		wg.Done()
	}()
	close(done)
	wg.Wait()
}

func main() {
	fmt.Println("Started context broadcast:")
	context_broadcast()
	fmt.Println("Started done channel:")
	done_channel()
}

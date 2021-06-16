package main

import (
	"context"
	"fmt"
	"sync"
)

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
	gouroutine_count := 2
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
	for i := 0; i < gouroutine_count; i++ { // I need to know how many there are
		done <- 1
	}
	wg.Wait()
}

func main() {
	fmt.Println("Started context broadcast:")
	// Context is better, it automatically broadcasts (by closing the Done() channel) to all goroutines that they should stop.
	context_broadcast()
	fmt.Println("Started done channel:")
	// To cancel all child goroutines with the done channel, I need to know how many gouroutines there are.
	done_channel()
}

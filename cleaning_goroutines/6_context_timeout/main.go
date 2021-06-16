package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Small example showing how context.WithTimeout can stop a function.

A context timeout has no power to stop a goroutine or a piece of running logic.
It is the responsibility of the goroutine to check if the timeout happened.

So, goroutines that can be timedout or cancelled should be broken into working bits, so that they can
check wether the context cancellation has occurred or not.
*/

func main() {
	wg := sync.WaitGroup{} // just to wait the program to finish.
	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// The context timeout timer starts running as soon as it is created!
	// Not as soon as it is passed to another function.
	time.Sleep(3 * time.Second) // only 1 "Doing my job" will be printed to illustrate this.

	defer cancel()
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-time.After(1 * time.Second):
				fmt.Println("Doing my job")
			case <-ctx.Done():
				fmt.Println("Context done")
				return
			}
		}
	}(ctx)

	<-ctx.Done()
	wg.Wait()
}

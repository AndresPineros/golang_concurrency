package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// When the done channel is triggered, all stages are cancelled

	generator := func(ctx context.Context) <-chan int {
		intStream := make(chan int, 30) // Buffered to make stages independent
		count := 0
		go func() {
			defer close(intStream)
			for {
				select {
				case <-time.After(200 * time.Millisecond):
					count++
					intStream <- int(rand.Uint32())
				case <-ctx.Done():
					fmt.Println("Cancelled generator stage. Produced", count, "integers")
					return
				}
			}
		}()
		return intStream
	}

	multiply := func(ctx context.Context, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int, 30) // Buffered to make stages independent
		count := 0
		go func() {
			defer close(multipliedStream)
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Source stream closed for multiply stage. Processed", count, "numbers and exited.")
					return
				case i := <-intStream:
					time.Sleep(300 * time.Millisecond)
					count++
					multipliedStream <- i * multiplier
				}
			}
		}()
		return multipliedStream
	}

	add := func(ctx context.Context, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int, 30) // Buffered to make stages independent
		count := 0
		go func() {
			defer close(addedStream)
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Source stream closed for add stage. Processed", count, "numbers and exited.")
					return
				case i := <-intStream:
					time.Sleep(500 * time.Millisecond)
					count++
					addedStream <- i + additive
				}
			}
		}()
		return addedStream
	}

	ctx, cancel := context.WithCancel(context.Background())
	// The context cancellation will stop all pipelines
	go func() {
		time.Sleep(5 * time.Second)
		cancel() // Stop the pipeline after 5 seconds
	}()
	intStream := generator(ctx)
	pipeline := multiply(ctx, add(ctx, multiply(ctx, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
	time.Sleep(5 * time.Second) // ugly, but just to wait for stages to print the count
}

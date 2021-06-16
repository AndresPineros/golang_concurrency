package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

/*
The workers from the previous channel don't write the results
to a channel. What if we need to pass these results to another goroutine via a channel?
*/

func produce(ctx context.Context, workerCount int) <-chan int {
	channel := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(workerCount)
	go func() {
		wg.Wait()
		close(channel)
	}()
	for i := 0; i < workerCount; i++ {
		go func(producerId int) {
			jobCounter := 0
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(1 * time.Second):
					select {
					case <-ctx.Done():
						return
					case channel <- jobCounter:
						jobCounter++
					}
				}
			}
		}(i)
	}
	return channel
}

/*
Consume reads work from a channel and writes results to another channel.
This is kinda like a stage of a pipeline with multiple workers.
*/

func consume(ctx context.Context, channel <-chan int, workerCount int) <-chan int {
	results := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)
	go func() {
		defer close(results)
		wg.Wait()
	}()
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job := <-channel:
					select {
					case <-ctx.Done():
						return
					case results <- job:
						time.Sleep(2 * time.Second)
					}
				}
			}
		}()
	}
	return results
}

func main() {
	writerCount := 3
	readerCount := 2

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	channel := produce(ctx, writerCount)
	results := consume(ctx, channel, readerCount)

	go func() {
		<-signals
		fmt.Println("Received interrupt. Cleaning resources.")
		cancel()
	}()

	for i := range results {
		fmt.Println("Result:", i)
	}
}

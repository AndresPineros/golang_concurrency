package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

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
				/*
					This rewrite shows the correct way of dealing with
					blocking actions inside a case when other cases must be listening.

					Also, check the or-done pattern.
				*/
				done := func() {
					fmt.Println("Producer", producerId, "done.", "Produced", jobCounter)
				}
				select {
				case <-ctx.Done():
					done()
					return
				case <-time.After(2 * time.Second):
					fmt.Println("Producer", producerId, "produced job", jobCounter)
					select {
					case <-ctx.Done():
						done()
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

func consume(channel <-chan int) {
	for job := range channel {
		fmt.Println("Worker received job", job)
		time.Sleep(3 * time.Second)
	}
	fmt.Println("Worker done")
}

func main() {
	writerCount := 3

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	channel := produce(ctx, writerCount)
	go consume(channel)

	<-signals
	fmt.Println("Received interrupt. Cleaning resources.")
	cancel()
}

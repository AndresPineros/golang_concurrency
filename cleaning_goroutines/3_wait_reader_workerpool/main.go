package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func produce(ctx context.Context) <-chan int {
	channel := make(chan int)
	jobCounter := 0
	go func() {
		defer close(channel)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Producer done. Produced", jobCounter)
				return
			case <-time.After(2 * time.Second):
				channel <- jobCounter
				jobCounter++
			}
		}
	}()
	return channel
}

func consume(channel <-chan int, wg *sync.WaitGroup, consumerId int) {
	wg.Add(1)
	defer wg.Done()
	for job := range channel {
		fmt.Println("Worker", consumerId, "received job", job)
		time.Sleep(3 * time.Second)
	}
	fmt.Println("Worker done", consumerId)
}

func main() {
	workerCount := 5
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	channel := produce(ctx)
	for i := 0; i < workerCount; i++ {
		go consume(channel, &wg, i)
	}

	<-signals
	fmt.Println("Received interrupt. Cleaning resources.")
	cancel()
	wg.Wait()
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Producer struct {
	Jobchan chan int
}

func (p *Producer) produce(ctx context.Context) {
	producing := true
	jobCounter := 0
	for producing {
		select {
		case <-ctx.Done():
			producing = false
		case <-time.After(2 * time.Second):
			p.Jobchan <- jobCounter
			jobCounter++
		}
	}
	fmt.Println("Producer done")
}

type Worker struct {
	Jobchan chan int
	wg      *sync.WaitGroup
	Id      int
}

func (w *Worker) consume(ctx context.Context) {
	w.wg.Add(1)
	consuming := true
	for consuming {
		select {
		case <-ctx.Done():
			consuming = false
			w.wg.Done()
		case job := <-w.Jobchan:
			fmt.Println("Worker", w.Id, "received job", job)
			for i := 0; i < 3; i++ {
				time.Sleep(1 * time.Second)
			}
		}
	}
	fmt.Println("Worker done", w.Id)
}

func main() {
	workerCount := 5
	jobchan := make(chan int)
	wg := &sync.WaitGroup{}
	producer := &Producer{Jobchan: jobchan}

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go producer.produce(ctx)
	for i := 0; i < workerCount; i++ {
		w := &Worker{Id: i, Jobchan: jobchan, wg: wg}
		go w.consume(ctx)
	}

	<-signals
	fmt.Println("Received interrupt. Cleaning resources.")
	cancel()
	wg.Wait()
}

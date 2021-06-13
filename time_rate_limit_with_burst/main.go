package main

import (
	"fmt"
	"sync"
	"time"
)

func hitAPI(opId int, wg *sync.WaitGroup, burstChan <-chan interface{}, done <-chan interface{}) {
	<-burstChan
	// Expensive operation that hits an API hard
	fmt.Println("Started operation with ID", opId)
	time.Sleep(15 * time.Second)
	fmt.Println("Finished operation with ID", opId)
	wg.Done()
}

func main() {
	fmt.Println("Start")
	burstCapacity := 10
	wg := sync.WaitGroup{} // just to wait all operations to finish and see the ratelimit being applied.
	burstChan := make(chan interface{}, burstCapacity)
	done := make(chan interface{})
	go func() {
		// initialize burst capacity
		for i := 0; i < burstCapacity; i++ {
			burstChan <- 1
		}
		tick := time.Tick(1 * time.Second)
		for {
			select {
			case burstChan <- <-tick: // If tick can write to burstchan.
				fmt.Println("Current burst capacity", len(burstChan))
			case <-done:
				return
			}
		}
	}()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go hitAPI(i, &wg, burstChan, done)
	}
	wg.Wait()
	done <- 1
}

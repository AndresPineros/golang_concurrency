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
			/*
				Doing case ch1<- <-ch2 is dangerous because
				if nobody is listening on ch1 the output of <-ch2 won't
				be written into ch1 and will be lost.

				It doesn't matter in this case because ticks can be lost.
			*/
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

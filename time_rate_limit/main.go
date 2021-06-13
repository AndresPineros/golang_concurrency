package main

import (
	"fmt"
	"sync"
	"time"
)

func hitAPI(opId int, wg *sync.WaitGroup, ratelimit <-chan time.Time) {
	<-ratelimit
	// Expensive operation that hits an API hard
	time.Sleep(3 * time.Second)
	fmt.Println("Finished operation with ID", opId)
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{}                  // just to wait all operations to finish and see the ratelimit being applied.
	ratelimit := time.Tick(1 * time.Second) // ticks every second, allowing a new operation to be executed

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go hitAPI(i, &wg, ratelimit)
	}
	wg.Wait()
}

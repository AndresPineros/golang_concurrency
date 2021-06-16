package main

import (
	"fmt"
	"time"
)

/*
Workerpools are a nice pattern to distribute work across several goroutines BUT
they are really bad because:

1) We're spawning goroutines that we don't know are going to be used. Even if they are cheap, they cost.
2) Having all these workers alive means that debugging the code is really hard. Whenever there is a deadlock
   it is normal to do a goroutine dump to see which goroutine is stuck. If we have a workerpool with 1000
   goroutines waiting for input, it will be extremely hard to know which one is the culprit.

The alternative is: Use semaphores to spawn goroutines ONLY when they are needed.
*/

func workerpool() {
	workers := 3
	workchan := make(chan int)
	for i := 0; i < workers; i++ {
		go func() {
			for i := range workchan {
				time.Sleep(3 * time.Second)
				fmt.Println("Workerpool worked on ", i)
			}
		}()
	}
	amountOfWork := 10
	for i := 0; i < amountOfWork; i++ {
		workchan <- i
	}
	fmt.Println("Finished workerpool work")
	//Give some time for goroutines to finish. To avoid using WaitGroup and loosing focus.
	time.Sleep(5 * time.Second)
}

type token struct{}

func semaphore() {
	workers := 3
	sem := make(chan token, workers)
	amountOfWork := 10
	for i := 0; i < amountOfWork; i++ {
		sem <- token{} // This is a blocking operation, so this loop will only exit when all tasks have been assined to a worker
		go func(task int) {
			time.Sleep(3 * time.Second)
			fmt.Println("Semaphore worked on ", task)
			<-sem
		}(i)
	}

	for i := 0; i < workers; i++ {
		sem <- token{}
	}
	fmt.Println("Finished semaphore work")
	//Give some time for goroutines to finish. To avoid using WaitGroup and loosing focus.
	time.Sleep(5 * time.Second)
}

func main() {
	// fmt.Println("Start workerpool example:")
	// workerpool()
	fmt.Println("Start semaphore example:")
	semaphore()
}

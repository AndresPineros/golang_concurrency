package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
This is a bad example. I feel the Wait() should go in the workers
and the Broadcast() in the control goroutine, not the other way around.
Anyways, it illustrates the workflow of the tool
*/
func main() {
	rand.Seed(time.Now().UnixNano())

	const N = 10
	var values [N]string

	cond := sync.NewCond(&sync.Mutex{})

	for i := 0; i < N; i++ {
		go func(i int) {
			fmt.Println(i, "Locked")
			cond.L.Lock()
			fmt.Println("Modifying data", i)
			time.Sleep(1 * time.Second)
			values[i] = fmt.Sprint('a' + i)
			cond.Broadcast()
			fmt.Println(i, "Unlocked")
			cond.L.Unlock()
		}(i)
	}

	// This function must be called when
	// cond.L is locked.
	checkCondition := func() bool {
		fmt.Println(values)
		for i := 0; i < N; i++ {
			if values[i] == "" {
				return false
			}
		}
		return true
	}

	time.Sleep(2 * time.Second)
	cond.L.Lock() // Lock to ensure wait is called while locked
	defer cond.L.Unlock()
	for !checkCondition() {
		fmt.Println("Waiting...")
		cond.Wait()
		fmt.Println("Waiting over...")
	}
}

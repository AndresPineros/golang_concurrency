package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Do readers receive a message first if they arrive first to ask for it?
Let's say I have a channel with 2 reader goroutines: r1 and r2.
If r1 asked to read from a channel before r2, will r1 get the message first?

Answer: This test can't prove it but it seems like it.
*/

func main() {

	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(readerId int) {
			<-ch
			fmt.Println("Reader", readerId, "received message")
			wg.Done()
		}(i)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(3 * time.Second) // to make sure all readers are requesting a message from the channel before we write.

	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		ch <- 0
	}
	wg.Wait()
}

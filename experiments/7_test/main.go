package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Are writers to channels FIFO just like readers?
Answer: This test can't prove it but it seems like it.
*/

func main() {

	ch := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(writerId int) {
			ch <- fmt.Sprintf("%v %v", "Writer", writerId)
			wg.Done()
		}(i)
		time.Sleep(100 * time.Millisecond)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := range ch {
		fmt.Println(i)
	}
}

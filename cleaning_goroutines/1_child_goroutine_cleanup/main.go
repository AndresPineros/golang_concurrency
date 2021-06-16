package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings: // nil channel will make case be ignored.
					fmt.Println(s)
				case <-done:
					fmt.Println("Done channel closed")
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	terminated := doWork(done, nil) // nil channel
	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}

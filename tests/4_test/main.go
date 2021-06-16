package main

import (
	"fmt"
	"time"
)

/*
What happens if a select tries to read from a nil channel?

Answer: reading and writing to a nil channel permanently blocks the goroutine.
	    but inside a select it does nothing
*/

func main() {
	nilchan := make(chan int)
	nilchan = nil
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Sleeping...")
		case v := <-nilchan:
			fmt.Println("Read from nilchan", v)
		case nilchan <- 100:
			fmt.Println("Write to nilchan")
		}
	}
}

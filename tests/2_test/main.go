package main

import (
	"fmt"
	"time"
)

/*
What happens first when redirecting from a channel to another in a select statement (like case ch1 <- <- ch2)

A) Does it receive first from channel 2 and then tries to send to channel 1?
B) Does it verify if it is possible to send to channel 1 before reading from channel 2?

Importance: If we read from ch2 without sending to ch1, the value that was read from ch2 will be lost!!

Answer:

*/

func main() {

	source := make(chan int, 0)
	results := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			source <- i
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			r := <-results
			fmt.Println("Received", r)
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				fmt.Println("Sleep")
			case results <- <-source:
				fmt.Println("Read")
			}
		}
	}()

	time.Sleep(50 * time.Second) // to avoid adding waitgroups and other blocks that distract from the main idea.
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func printStuff(done chan interface{}) {
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Working...")
		case <-done:
			fmt.Println("Stopped working...")
			return
		}
	}
}

func main() {
	done := make(chan interface{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		done <- <-c
	}()

	printStuff(done)
}

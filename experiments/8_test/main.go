package main

import (
	"fmt"
	"time"
)

/*
Closing a channel triggers a select statement case?
*/

func main() {
	ch := make(chan interface{})

	go func() {
		defer close(ch)
		time.Sleep(3 * time.Second)
	}()
	select {
	case <-ch:
		fmt.Println("Channel closed")
	}
}

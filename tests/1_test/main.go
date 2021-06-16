package main

import (
	"fmt"
	"time"
)

/*
What happens first when writing on a select statement?

A) Is the operation evaluated OR
B) is the receiving channel availability verified?

Answer: it first executes the operation and THEN tries to send it.
If there is no receiver on the other side of the channel, it moves on... BUT
the operation did ocurr, and if the operation mutates an object, the object will
be mutated even if it isn't sent.
*/

func main() {

	var value int = 5
	p := &value

	process := func(param *int) int {
		*param++
		return *param
	}

	results := make(chan int)

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
				fmt.Println("Sleep", *p)
			case results <- process(p): // process happens even if nobody is listening at results.
				fmt.Println("Processed", *p)
			}
		}
	}()

	time.Sleep(50 * time.Second) // to avoid adding waitgroups and other blocks that distract from the main idea.
}

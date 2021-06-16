package main

import "fmt"

/*
Are channels pointers or values?

Answer: they are pointers
*/

func main() {
	ch := make(chan string, 1)

	x := ch

	x <- "They are pointers"

	fmt.Println(<-ch)
}

package main

/*
Can I close a read-only channel?

Answer: No. close(producer()) won't compile.
*/

func producer() <-chan int {
	c := make(chan int)
	//...
	return c
}

func main() {
	//close(producer())
}

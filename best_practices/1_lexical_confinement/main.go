package main

import "fmt"

func main() {
	fmt.Println("Running lexical confinement:")
	lexical_confinement()
	fmt.Println("Running adhoc confinement:")
	ad_hoc_confinement()

	/*
		Both functions do the same BUT the lexical confinement hides the writable channel from the outter world.
		This means that nobody can mess with the channel (close it, write to it other stuff).

		The ad_hoc confinment means that we expect our team to understand that they shouldn't do this, but
		they can forget and make mistakes.

		Lexical confinement is better because it prevents mistakes.
	*/
}

func lexical_confinement() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5) // Only the chanOwner function can write or close this channel.
		go func() {
			defer close(results)
			data := []int{0, 1, 2, 3, 4, 5}
			for i := range data {
				results <- i
			}
		}()
		return results
	}
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	results := chanOwner()
	consumer(results)
}

func ad_hoc_confinement() {
	data := []int{0, 1, 2, 3, 4, 5} // This data could be manipulated by anyone in this scope.
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int) // This channel could be manipulated by anyone in this scope.
	go loopData(handleData)
	for num := range handleData {
		fmt.Println("Received:", num)
	}
	fmt.Println("Done receiving!")
}

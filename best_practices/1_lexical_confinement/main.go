package main

import "fmt"

func main() {
	fmt.Println("Running lexical confinement:")
	lexical_confinement()
	fmt.Println("Running adhoc confinement:")
	ad_hoc_confinement()

	/*
		Both functions do the same. The difference is that in the ad_hoc confinement, the data is
		confined but a developer could break the confinement and access the data from multiple goroutines.

		In the lexical confinement, it is impossible to break the confinment. The chanOwner function is the only place
		where data can be written to the channel because it returns a receive-only channel (the <-chan int).
		Also, it is the only entity in charge of closing the channel.
	*/
}

func lexical_confinement() {
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
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
	data := []int{0, 1, 2, 3, 4, 5} // This data could be accessed concurrently if a developer changes the code.

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)
	for num := range handleData {
		fmt.Println("Received:", num)
	}
	fmt.Println("Done receiving!")
}

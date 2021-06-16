package main

import (
	"fmt"
	"time"
)

/*
Didn't like the way those selects look, so I rewrote it.

There's something I don't like about the heartbeat pattern. It uses channels and it means it can only have
a single client, because with channels messages aren't broadcasted to all readers.
*/

func doWork(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		//helpers
		workGenTmp := workGen
		var resultsTmp chan time.Time
		var result time.Time
		//helpers

		for {
			select {
			case <-done:
				return
			case <-pulse:
				select {
				case heartbeat <- struct{}{}:
				default: // This deault helps in case nobody is reading the heartbeat.
				}
			case result = <-workGenTmp: // If we're receiving work, we enable the send work case and disable receiving work case.
				workGenTmp = nil
				resultsTmp = results
			case resultsTmp <- result: // If we're sending work, we enable the receive work case and disable the send work case.
				workGenTmp = workGen
				resultsTmp = nil
			}
			/*
				This means we're either sending or receiving work, not both at the same time, and we're not blocking the heartbeat pulse.
			*/
		}
	}()
	return heartbeat, results
}

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })
	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}

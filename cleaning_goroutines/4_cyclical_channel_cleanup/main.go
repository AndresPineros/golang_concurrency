/*
This is a little rewrite of https://golang.org/doc/codewalk/sharemem/

I wanted to have a circular channel in which the writers of a channel are also the readers.
How do we close that channel safely to avoid writing to a closed channel (panic)?

I found a way to do it with

Also, the writers/readers are multiple workers, not a single one.
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	numPollers     = 2                // number of Poller goroutines to launch
	pollInterval   = 2 * time.Second  // how often to poll each URL
	statusInterval = 10 * time.Second // how often to log status to stdout
	errTimeout     = 10 * time.Second // back-off timeout on error
)

var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
}

type Resource struct {
	url      string
	errCount int
}

func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

func Poller(safeWrite *bool, m *sync.RWMutex, resources chan *Resource) {
	for r := range resources {
		r.Poll()
		fmt.Println("Processing", r.url)
		go func(r *Resource) {
			time.Sleep(1 * time.Second)
			m.RLock()
			if *safeWrite {
				resources <- r
			}
			m.RUnlock()
		}(r)
	}
}

func main() {
	resources := make(chan *Resource) // This is the circular channel. The Pollers read from it and write back to it.
	ctx, cancel := context.WithCancel(context.Background())
	safeWrite := true
	mux := sync.RWMutex{}
	time.AfterFunc(20*time.Second, func() { cancel() }) // cancel program after 20 seconds.
	go func() {
		<-ctx.Done()
		mux.Lock()
		safeWrite = false
		close(resources)
		mux.Unlock()
	}()

	for i := 0; i < numPollers; i++ {
		go Poller(&safeWrite, &mux, resources)
	}
	for _, url := range urls {
		resources <- &Resource{url: url}
	}
	time.Sleep(10 * time.Second)
}

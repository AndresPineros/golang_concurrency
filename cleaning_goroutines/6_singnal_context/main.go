package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	select {
	case <-time.After(time.Second * 3):
		fmt.Println("missed signal")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		cancel()
	}

}

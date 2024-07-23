package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup) {
	wg.Add(1)

	fmt.Println("Worker start")
	time.Sleep(5 * time.Second)
	fmt.Println("Worker finish")

	wg.Done()
}

func cleanupServer() {
	fmt.Println("Server cleanup")
}

func main() {
	defer cleanupServer()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	wg := sync.WaitGroup{}

	for {
		time.Sleep(1 * time.Second)
		brk := false
		select {
		case <-ctx.Done():
			brk = true
		default:
		}
		if brk {
			break
		}
		go worker(&wg)
	}
	stop()
	fmt.Println("finishing active requests, press Ctrl+C again to force shutdown")
	wg.Wait()
}

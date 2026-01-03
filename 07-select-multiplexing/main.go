package main

import (
	"fmt"
	"sync"
	"time"
)

// Producer function that simulates a bot sending messages at intervals.
// It sends integers to its output channel and listens for a done signal to terminate early.
// When done is closed, the producer stops sending messages.

func producer(done chan struct{}, botName string, delay time.Duration, wg *sync.WaitGroup) chan int {
	out := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		for i := 1; i <= 5; i++ {
			time.Sleep(delay)
			select {
			case <-done:
				fmt.Println(botName, "received done signal")
				return
			case out <- i:
				fmt.Println(botName, "produced", i)
			}
		}
	}()
	return out
}
func main() {
	fmt.Println("Starting Select Multiplexing example")
	done := make(chan struct{})

	wg := sync.WaitGroup{}

	ch1 := producer(done, "Bot-1", 500*time.Millisecond, &wg)
	ch2 := producer(done, "Bot-2", 1000*time.Millisecond, &wg)

	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case val, ok := <-ch1:
			if !ok {
				ch1 = nil
			}
			if val != 0 {
				fmt.Println("Main received from Bot-1:", val)
			}
		case val, ok := <-ch2:
			if !ok {
				ch2 = nil
			}
			if val != 0 {
				fmt.Println("Main received from Bot-2:", val)
			}
		case <-time.After(800 * time.Millisecond):
			fmt.Println("Timeout! Algum bot demorou demais entre mensagens.")
			done <- struct{}{}
		}
		if ch1 == nil && ch2 == nil {
			break
		}

	}
	fmt.Println("All bots have finished. Exiting main.")

}

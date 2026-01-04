package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// worker simulates a worker that processes integers from the input channel,
// doubles them after a delay, and sends them to the output channel.
// It listens for context cancellation to terminate gracefully.
// The worker also handles the case where the input channel is closed.

func worker(ctx context.Context, id int, wg *sync.WaitGroup, in <-chan int, out chan<- int, delay time.Duration) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d received done signal\n", id)
			return
		case val, ok := <-in:
			if !ok {
				fmt.Printf("worker %d input channel closed\n", id)
				return
			}
			fmt.Printf("worker %d received %d\n", id, val)
			time.Sleep(delay)

			select {
			case <-ctx.Done():
				fmt.Printf("worker %d received done signal\n", id)
				return
			case out <- val * 2:
				fmt.Printf("worker %d processed %d to %d\n", id, val, val*2)
			}
		}
	}
}

func main() {

	numWorkers, numDispatchers := rand.Intn(5)+1, rand.Intn(5)+1
	fmt.Printf("Starting with %d workers and %d dispatchers\n", numWorkers, numDispatchers)
	var wgDispatchers sync.WaitGroup
	var wgWorkers sync.WaitGroup

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	in := make(chan int)
	out := make(chan int)

	for i := 1; i <= numDispatchers; i++ {
		wgDispatchers.Add(1)
		go func(id int) {
			defer wgDispatchers.Done()
			for j := 0; j < 3; j++ {
				val := rand.Intn(100)
				select {
				case <-ctx.Done():
					fmt.Printf("dispatcher %d received done signal\n", id)
					return
				case in <- val:
					fmt.Printf("dispatcher %d dispatched a value %d\n", id, val)
				}
			}
		}(i)
	}

	for i := 1; i <= numWorkers; i++ {
		wgWorkers.Add(1)
		go worker(ctx, i, &wgWorkers, in, out, time.Duration(500+rand.Intn(1500))*time.Millisecond)
	}

	go func() {
		wgDispatchers.Wait()
		close(in)
		wgWorkers.Wait()
		close(out)
	}()

	for result := range out {
		fmt.Println("Main received processed value:", result)
	}

	fmt.Println("All workers have finished. Exiting main.")

}

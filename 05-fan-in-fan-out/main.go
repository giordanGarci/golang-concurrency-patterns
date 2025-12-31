package main

import (
	"fmt"
	"sync"
	"time"
)

// Fan-out: We launch multiple Goroutines to perform work concurrently.
// Fan-in: We collect the results from multiple Goroutines into a single channel.
func producer(done chan struct{}) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= 10; i++ {
			select {
			case out <- i:
				fmt.Println("Producing", i)
				time.Sleep(100 * time.Millisecond) // Simulate work
			case <-done:
				fmt.Println("Producer received done signal")
				return
			}
		}
	}()
	return out
}

func processor(id int, in chan int, out chan int, wg *sync.WaitGroup, done chan struct{}) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Printf("Processor %d received done signal\n", id)
				return
			case i, ok := <-in:
				if !ok {
					return
				}
				result := i * 2
				select {

				case <-done:
					fmt.Printf("Processor %d received done signal\n", id)
					return
				case out <- result:
					fmt.Printf("Processor %d processing %d, result %d\n", id, i, result)
				}
			}
		}
	}()
}

func main() {
	fmt.Println("Starting Fan-in/Fan-out example")

	done := make(chan struct{})
	in := producer(done)
	out := make(chan int)

	// Fan-out: Start multiple processors
	numProcessors := 3
	wg := sync.WaitGroup{}

	for i := 1; i <= numProcessors; i++ {
		wg.Add(1)
		processor(i, in, out, &wg, done)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	// Fan-in: Collect results
	counter := 0
	for result := range out {
		counter++
		fmt.Println("Received result:", result)
		if counter == 5 {
			fmt.Println("Received 5 results, sending done signal")
			close(done)
		}
	}

}

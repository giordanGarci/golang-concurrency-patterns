package main

import (
	"fmt"
	"sync"
)

//  This example demonstrates a concurrent pipeline with a generator and a processor.
// The generator produces numbers, and the processor doubles them.
// The pipeline can be gracefully terminated using a done channel.

func generator(done chan struct{}) chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= 10; i++ {
			select {
			case out <- i:
				fmt.Println("Generating", i)
			case <-done:
				fmt.Println("Generator received done signal")
				return
			}
		}
	}()
	return out
}

func processor(in chan int, out chan int, done chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("Processor received done signal")
				return
			case i, ok := <-in:
				if !ok {
					return
				}
				result := i * 2

				select {

				case out <- result:
					fmt.Printf("Processor processing %d, result %d\n", i, result)
				case <-done:
					fmt.Println("Processor received done signal")
					return
				}
			}
		}
	}()
}

func main() {
	fmt.Println("Starting Concurrent Pipeline example")
	done := make(chan struct{})
	in := generator(done)
	out := make(chan int)

	wg := sync.WaitGroup{}
	processor(in, out, done, &wg)

	go func() {
		wg.Wait()
		close(out)
	}()

	for result := range out {
		fmt.Println("Final result:", result)
	}

}

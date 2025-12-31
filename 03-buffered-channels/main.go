package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(3 * time.Second)
		// Simulate doing some work
		fmt.Printf("Worker %d processing job %d\n", id, j)
		results <- j * 2
	}
}

func main() {
	const numWorkers = 3

	jobs := make(chan int, 10)
	results := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Send jobs
	go func() {
		for j := 1; j <= 9; j++ {
			fmt.Println("Sending job", j)
			jobs <- j
		}
		close(jobs)
	}()

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, &wg)
	}

	go func() {
		wg.Wait()      // when all workers are done
		close(results) // ... close results channel
	}()

	for res := range results {
		fmt.Println("Result received:", res)
	}

	fmt.Println("All jobs processed.")
}

package main

import (
	"fmt"
	"sync"
)

// 1000 Goroutines (go func()), where all of them increment the same variable concurrently.
// What to observe: If we use a mutex to protect the increment operation, the final result will always be 1000.
// This shows that the mutex successfully
func main() {

	var contador int
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 1; i <= 1000; i++ {
		go func() {
			mu.Lock()
			contador++
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Contador final:", contador)
}

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker simulate a task that takes a random amount of time to complete.
// It respects the context's timeout and cancels its work if the context is done.

func worker(id int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	slowJob := time.Duration(2+rand.Intn(7)) * time.Second

	fmt.Printf("Worker %d: Preciso de %v para terminar...\n", id, slowJob)

	select {
	case <-time.After(slowJob):
		fmt.Printf("Worker %d: Terminei meu trabalho!\n", id)
	case <-ctx.Done():
		fmt.Printf("Worker %d: Tempo esgotado! Cancelando trabalho...\n", id)
	}

}

func main() {
	fmt.Println("Iniciando exemplo de Context com Timeout")
	numWorkers := 6
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, ctx, &wg)
	}

	wg.Wait()
	fmt.Println("Todos os workers finalizaram. Saindo do main.")

}

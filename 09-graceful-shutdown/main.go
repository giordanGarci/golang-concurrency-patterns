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

// this program demonstrates a graceful shutdown mechanism using context and OS signals.
// It spawns multiple worker goroutines that perform tasks and listen for termination signals.
// Upon receiving a signal, the main function notifies all workers to finish their tasks gracefully.

func main() {
	fmt.Println("Iniciando exemplo de Graceful Shutdown")
	numWorkers := 6
	var wg sync.WaitGroup

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			randTime := time.Duration(2+rand.Intn(7)) * time.Second
			fmt.Printf("Worker %d: Iniciando trabalho...\n", id)

			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d: Recebi sinal de término! Finalizando trabalho...\n", id)
				return
			case <-time.After(randTime):
				fmt.Printf("Worker %d: Trabalho concluído!\n", id)
			}

		}(i)
	}

	<-ctx.Done()
	fmt.Println("Main: Recebi sinal de término! Aguardando workers finalizarem...")

	wg.Wait()
	fmt.Println("Todos os workers finalizaram. Saindo do main.")
}

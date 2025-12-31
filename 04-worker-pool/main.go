package main

import (
	"fmt"
)

func generator(done chan struct{}) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= 10; i++ {
			select {
			case <-done:
				fmt.Println("Generator received done signal")
				return
			case out <- i:
				fmt.Println("Generated:", i)
			}
		}
	}()
	return out
}

func transformer(done chan struct{}, in chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-done:
				fmt.Println("Transformer received done signal")
				return
			case out <- n * 2:
				fmt.Println("Transformed:", n, "to", n*2)
			}
		}
	}()
	return out
}

func saver(done chan struct{}, in chan int) {
	var count int
	for n := range in {
		count++
		fmt.Println("Saved:", n)
		if count == 5 {
			fmt.Println("Saver reached limit, sending done signal")
			close(done)
			return
		}
	}
}

func main() {
	done := make(chan struct{})

	saver(done, transformer(done, generator(done)))
}

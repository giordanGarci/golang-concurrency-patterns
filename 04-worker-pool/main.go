package main

import (
	"fmt"
)

func generator() chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			out <- i
			fmt.Println("Generated:", i)
		}
		close(out)
	}()
	return out
}

func transformer(in chan int) chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * 2
			fmt.Println("Transformed:", n, "to", n*2)
		}
		close(out)
	}()
	return out
}

func saver(in chan int) {
	for n := range in {
		fmt.Println("Saved value:", n)
	}
}

func main() {
	saver(transformer(generator()))
}

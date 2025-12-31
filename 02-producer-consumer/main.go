package main

import (
	"fmt"
	"time"
)

// This example demonstrates a simple producer-consumer pattern using a buffered channel.
// The producer generates integers and sends them to the channel.
// The consumer receives integers from the channel and processes them with a delay.

func producer(c chan<- int) {
	for i := range 10 {
		fmt.Println("Producing:", i)
		c <- i
	}
	close(c)
}

func main() {
	c := make(chan int, 1)
	go producer(c)

	for v := range c {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Received:", v)
	}

}

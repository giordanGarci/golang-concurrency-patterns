package main

import (
	"fmt"
	"time"
)

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

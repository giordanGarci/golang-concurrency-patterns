package main

import (
	"fmt"
)

// 1000 Goroutines (go func()), where all of them increment the same variable concurrently.
// What to observe: In the end, the result will rarely be 1000. It will be an odd number (942, 875, etc.).
// This proves that some operations "overlapped" others.

// Why does it happen?
// The increment (i++) seems like a single operation, but for the processor it's actually three steps:
// Read the current value from memory.
// Add 1 to that value.
// Write the new value back to memory.
// If Goroutine A reads the value (say, 10) and, before it writes 11, Goroutine B also reads the same 10, both will write 11. The result should have been 12, but one increment was "lost".
func increment(contador *int) {
	*contador++
}
func main() {
	contador := 0
	for i := 1; i <= 1000; i++ {
		go increment(&contador)
	}

	fmt.Println("Contador final:", contador)
}

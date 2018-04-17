package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup // initialize the WaitGroup

// This program should go to 11, but sometimes it only prints 1 to 10.
func main() {
	ch := make(chan int)
	wg.Add(1) // Add to the waitgroup, making sure the print is
	go Print(ch)
	for i := 1; i <= 11; i++ {
		ch <- i
	}
	close(ch) // Close the channel
	wg.Wait() // Wait for all the routines to be marked as "Done"
}

// Print prints all numbers sent on the channel.
// The function returns when the channel is closed.
func Print(ch <-chan int) {
	for n := range ch { // reads from channel until it's closed
		time.Sleep(1000) // For easier noticing the bug
		fmt.Println(n)
	}
	wg.Done() // Mark routine as "Done"
}

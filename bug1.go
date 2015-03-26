package main

import (
	"fmt"
)

// I want this program to print "Hello world!", but it doesn't work.
func main() {
	ch := make(chan string)
	go func() {
		ch <- "Hello world!"
		close(ch)
	}()
	fmt.Println(<-ch)
}

/* All goroutines are asleep - deadlock!
asd


*/

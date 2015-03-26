package main

import "fmt"

// I want this program to print "Hello world!", but it doesn't work.
func main() {
	ch := make(chan string)
	go func() {
		ch <- "Hello world!"
	}()
	fmt.Println(<-ch)
}

/* All goroutines are asleep - deadlock!
In order to send to a channel, it needs to be within a go routine. I simply made a go func() to do the sending.


*/

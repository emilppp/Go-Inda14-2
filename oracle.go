// Stefan Nilsson 2013-03-13

// This program implements an ELIZA-like oracle (en.wikipedia.org/wiki/ELIZA).
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	star   = "Pythia"
	venue  = "Delphi"
	prompt = "> "
)

func main() {
	fmt.Printf("Welcome to %s, the oracle at %s.\n", star, venue)
	fmt.Println("Your questions will be answered in due time.")

	oracle := Oracle()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Printf("%s heard: %s\n", star, line)
		oracle <- line // The channel doesn't block.
	}
}

// Oracle returns a channel on which you can send your questions to the oracle.
// You may send as many questions as you like on this channel, it never blocks.
// The answers arrive on stdout, but only when the oracle so decides.
// The oracle also prints sporadic prophecies to stdout even without being asked.
func Oracle() chan<- string {
	questions := make(chan string)
	answers := make(chan string)
	go func() { // Send the question from questions to the prophecy method
		for {
			s := <-questions
			go func() {
				prophecy(s, answers)
			}()
		}
	}()

	go func() { //Wait some time and then give a random prophecy
		for {
			time.Sleep(time.Duration(5+rand.Intn(8)) * time.Second)
			prophecy("", answers)
		}
	}()

	go func() { //Makes the output print out, letter after letter, with 50 ms delay, making it look like the oracle is writing it.
		for {
			answer := <-answers
			answerIArray := strings.Split(answer, "")
			fmt.Println()
			for _, bokstav := range answerIArray {
				time.Sleep(time.Millisecond * 50)
				fmt.Print(bokstav)
			}
			fmt.Println()
			fmt.Print(prompt)
		}
	}()
	// TODO: Answer questions. X
	// TODO: Make prophecies. X
	// TODO: Print answers. X
	return questions
}

// This is the oracle's secret algorithm.
// It waits for a while and then sends a message on the answer channel.
// TODO: make it better.
func prophecy(question string, answer chan<- string) {
	// Keep them waiting. Pythia, the original oracle at Delphi,
	// only gave prophecies on the seventh day of each month.
	time.Sleep(time.Duration(20+rand.Intn(10)) * time.Second)

	// Find the longest word.
	longestWord := ""
	words := strings.Fields(question) // Fields extracts the words into a slice.
	for _, w := range words {
		if len(w) > len(longestWord) {
			longestWord = w
		}
	}

	// Cook up some pointless nonsense.
	nonsense := []string{
		"The moon is dark.",
		"The sun is bright.",
		"The water is blött.",
		"I can't be arsed.",
		"Would you mind?",
		"Can't you tell im busy?",
	}
	answer <- longestWord + " ...sigh, " + nonsense[rand.Intn(len(nonsense))]
}

func init() { // Functions called "init" are executed before the main function.
	// Use new pseudo random numbers every time.
	rand.Seed(time.Now().Unix())
}

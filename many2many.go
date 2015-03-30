// Stefan Nilsson 2013-03-13

// This is a testbed to help you understand channels better.
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Use different random numbers each time this program is executed.
	rand.Seed(time.Now().Unix())

	const strings = 100
	const producers = 4
	const consumers = 4

	before := time.Now()
	ch := make(chan string)
	wgp := new(sync.WaitGroup)
	wgp2 := new(sync.WaitGroup)
	wgp.Add(producers) // Make two waitgroups, one for producers and one for consumers
	wgp2.Add(consumers)
	for i := 0; i < producers; i++ {
		go Produce("p"+strconv.Itoa(i), strings/producers, ch, wgp)
	}
	for i := 0; i < consumers; i++ {
		go Consume("c"+strconv.Itoa(i), ch, wgp2)
	}
	wgp.Wait() // Wait for all producers to finish.
	close(ch)  // Close the channel and then wait for the consumers to be done before printing time.
	wgp2.Wait()
	fmt.Println("time:", time.Now().Sub(before))
	select {}
}

// Produce sends n different strings on the channel and notifies wg when done.
func Produce(id string, n int, ch chan<- string, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		RandomSleep(100) // Simulate time to produce data.
		ch <- id + ":" + strconv.Itoa(i)
	}
	wg.Done() // Alert main that producers are done.
}

// Consume prints strings received from the channel until the channel is closed.
func Consume(id string, ch <-chan string, wg *sync.WaitGroup) { // Consume now takes a waitgroup aswell.
	for s := range ch {
		fmt.Println(id, "received", s)
		RandomSleep(100) // Simulate time to consume data.
	}
	wg.Done() // Alert main that consumers are done.
}

// RandomSleep waits for x ms, where x is a random number, 0 â‰¤ x < n,
// and then returns.
func RandomSleep(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Millisecond)
}

/*
Förklara vad som händer och varför det händer om man gör följande ändringar i programmet.
Prova att först tänka ut vad som händer och testa sedan din hypotes genom att ändra och köra programmet.


Vad händer om man byter plats på satserna wgp.Wait() och close(ch) i slutet av main-funktionen?
- Hypotes: Kanalen stängs innan alla producenter är klara.
- Test: 'panic: send on closed channel'. Den vill skicka saker till en kanal som är stängd!
Vad händer om man flyttar close(ch) från main-funktionen och i stället stänger kanalen i slutet av funktionen Produce?
- Hypotes: Produce kanske hinner stänga kanalen innan consume har hunnit fått sista värdet
- Test: 'panic: send on closed channel'. Kanalen hann stängas i slutet!
Vad händer om man tar bort satsen close(ch) helt och hållet?
- Hypotes: Kanalen stängs inte. Kommer inte att hända något annat än det.
- Test: Funkade som innan. Kanalen skickas till consume där den tar ut värdena i kanalen.
Vad händer om man ökar antalet konsumenter från 2 till 4?
- Hypotes: Den kommer att köras med 4 trådar istället för 2.
- Test: Det gick mycket snabbare! Ju fler consumers desto snabbare går det.
Kan man vara säker på att alla strängar blir utskrivna innan programmet stannar?
- Hypotes: Om man också sätter ett wgp.Wait() resp. wg.Done() för consume. Det finns redan en för Produce som gör att
den kommer att köras tills alla strängar är skapade och skickade in till kanalen. Med en liknande grej för consume
borde det gå att garantera.
- Test: La till consumers i wgp och ändrade consumers metodhuvud för att skicka med wgp'n. Nu borde det funka som det gjorde innan
både för producers och consumers
*/

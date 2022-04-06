package main

import (
	"math/rand"
	"time"
)

func generate(from population, to chan<- [2]*patient) {
	// remove patients from pupulation (make method which is thread safe)
	// which cant interact anymore.
}

// simulateInteraction
func simulateInteraction(to <-chan [2]*patient) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C: // dormant for to long
				ticker.Stop()
				return

			case intercourse := <-to:
				go intercourse[0].interact(intercourse[1])

			}
			ticker.Reset(10 * time.Second)
		}
	}()
}

func main() {
	// parse flags / args

	rand.Seed(time.Now().Unix())

	// generate list from flags (len).
	// start generators to pick 2 random pacients from.
}

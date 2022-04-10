package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

var (
	popSizeFlag            = flag.Int("population-size", 10, "population-size determines the sample size")
	popInfectedFlag        = flag.Int("population-infected", 5, "population-infected determines the number of starting infected sample")
	popMaxInteractionsFlag = flag.Int("population-max-interactions", 3, "population-max-interactions determines the max number of interactions a pacient can have")
	infectionRateFlag      = flag.Float64("infection-rate", 0.5, "infection-rate determines the chance of disease transmition, ie: 0.5, 0.25, 1.0")

	processGeneratorsFlag = flag.Int("process-generators", 1, "process-generators determines the number of generators used")
	processProcessorsFlag = flag.Int("process-processors", 1, "process-processors determines the number of processors used")
)

// generate generates 2 random pacients from population field: from.
//
// safe to use concurrently as pacients are thread safe.
func generate(from *population, to chan<- [2]*patient) {
	for {
		// find suitable first pacient.
		pat1, key := from.Random()
		for !pat1.canInteract(nil) {
			from.RemoveKey(key)
			pat1, key = from.Random()

			if pat1 == nil { // no more patients.
				return
			}
		}

		// find suitable second pacient.
		pat2, key := from.Random()
		for !pat2.canInteract(nil) {
			from.RemoveKey(key)
			pat2, key = from.Random()
		}

		to <- [2]*patient{pat1, pat2}
	}
}

// simulateInteraction processes interactions from (from) with a 10 second timeout.
// run in own gorutine.
func simulateInteraction(from <-chan [2]*patient) {
	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C: // dormant for to long
			timer.Stop()
			return

		case intercourse := <-from:
			intercourse[0].interact(intercourse[1])
		}
		// after each process reset the timer.
		timer.Reset(10 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	pop := NewPopulation(
		*popSizeFlag,
		*popInfectedFlag,
		*popMaxInteractionsFlag,
		float32(*infectionRateFlag),
	)

	log.Println("BEFORE ------------------------")
	pop.Print()

	wait := &waiter{}
	// add generators.
	for i := 0; i < *processGeneratorsFlag; i++ {
		wait.addGenerator(generate)
	}
	// add processors.
	for i := 0; i < *processProcessorsFlag; i++ {
		wait.addProcessor(simulateInteraction)
	}

	// start simulation.
	wait.start(pop, make(chan [2]*patient))

	// wait for simulation to end.
	wait.wait()

	log.Println("AFTER ------------------------")
	pop.Print()
}

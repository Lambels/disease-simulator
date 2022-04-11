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

	processGeneratorsFlag       = flag.Int("process-generators", 1, "process-generators determines the number of generators used")
	processProcessorsFlag       = flag.Int("process-processors", 1, "process-processors determines the number of processors used")
	processProcessorTimeoutFlag = flag.Int("process-processor-timeout", 10, "process-processor-timeout determines the number of seconds used as processor timeout")
)

// generate generates 2 random pacients from population field: from.
//
// safe to use concurrently as pacients are thread safe.
func generate(from *population, to chan<- [2]*pacient) {
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

			if pat2 == nil { // no more patients.
				return
			}
		}

		to <- [2]*pacient{pat1, pat2}
	}

}

// simulateInteraction processes interactions from (from).
// run in own gorutine.
func simulateInteraction(from <-chan [2]*pacient) {
	timer := time.NewTimer(time.Duration(*processProcessorTimeoutFlag) * time.Second)
	for {
		select {
		case <-timer.C: // dormant for to long
			timer.Stop()
			return

		case intercourse := <-from:
			if err := intercourse[0].interact(intercourse[1]); err != nil {
				break
			}
			intercourse[1].interact(intercourse[0])
		}
		// after each process reset the timer.
		timer.Reset(time.Duration(*processProcessorTimeoutFlag) * time.Second)
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
	wait.start(pop, make(chan [2]*pacient))

	// wait for simulation to end.
	wait.wait()

	log.Println("AFTER ------------------------")
	pop.Print()
}

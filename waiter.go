package main

import (
	"log"
	"sync"
)

type waiter struct {
	wg sync.WaitGroup

	generators []func(*population, chan<- [2]*pacient)
	processors []func(<-chan [2]*pacient)
}

func (w *waiter) addGenerator(cmd func(*population, chan<- [2]*pacient)) {
	w.wg.Add(1)
	w.generators = append(w.generators, cmd)
}

func (w *waiter) addProcessor(cmd func(<-chan [2]*pacient)) {
	w.wg.Add(1)
	w.processors = append(w.processors, cmd)
}

func (w *waiter) start(pop *population, comChan chan [2]*pacient) {
	// start processors.
	for _, pro := range w.processors {
		go func(cmd func(<-chan [2]*pacient)) {
			cmd(comChan)
			w.wg.Done()
		}(pro)
	}
	log.Println("Processors Started.")

	// start generators.
	for _, gen := range w.generators {
		go func(cmd func(*population, chan<- [2]*pacient)) {
			cmd(pop, comChan)
			w.wg.Done()
		}(gen)
	}
	log.Println("Generators Started.")
}

func (w *waiter) wait() {
	w.wg.Wait()
}

package main

import "sync"

type waiter struct {
	wg *sync.WaitGroup

	generators []func(population, chan<- [2]*patient)
	processors []func(<-chan [2]*patient)
}

func (w *waiter) addGenerator(cmd func(population, chan<- [2]*patient)) {
	w.wg.Add(1)
	w.generators = append(w.generators, cmd)
}

func (w *waiter) addProcessor(cmd func(<-chan [2]*patient)) {
	w.wg.Add(1)
	w.processors = append(w.processors, cmd)
}

func (w *waiter) start(pop population, comChan chan [2]*patient) {
	// start generators.
	for _, gen := range w.generators {
		go func(cmd func(population, chan<- [2]*patient)) {
			cmd(pop, comChan)
			w.wg.Done()
		}(gen)
	}

	// start processors
	for _, pro := range w.processors {
		go func(cmd func(<-chan [2]*patient)) {
			cmd(comChan)
			w.wg.Done()
		}(pro)
	}
}

func (w *waiter) wait() {
	w.wg.Wait()
}

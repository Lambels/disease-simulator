package main

import (
	"math/rand"
	"sync"
)

type population struct {
	mu    sync.RWMutex
	clean map[int]*patient
	dirty map[int]*patient
}

func newPopulation(len, infected, maxInteractions int, rate float32) *population {
	if infected > len {
		infected = len
	}

	clean := make(map[int]*patient, len)
	dirty := make(map[int]*patient, len)
	for i := 0; i < len; i++ {
		pat := &patient{
			maxInteractions: maxInteractions,
			infectionChance: rate,
		}
		clean[i] = pat
		dirty[i] = pat
	}

	pop := &population{
		clean: clean,
		dirty: dirty,
	}

	// infect random patients.
	for i := 0; i < infected; i++ {
		pop.infectRandom()
	}

	return pop
}

// infectRandom recursivly infects random uninfected pacients.
func (p *population) infectRandom() {
	pos := rand.Intn(len(p.dirty))
	if pat := p.dirty[pos]; !pat.infected {
		pat.infected = true
		return
	} else {
		p.infectRandom()
	}
}

func (p *population) loadPopulation() map[int]*patient {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dirty
}

func (p *population) removeKey(key int) {
	p.mu.Lock()
	delete(p.dirty, key)
	p.mu.Unlock()
}

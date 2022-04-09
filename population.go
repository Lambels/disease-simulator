package main

import (
	"math/rand"
	"sync"
)

type population struct {
	mu       sync.RWMutex
	patients map[int]*patient
}

func newPopulation(len, infected, maxInteractions int, rate float32) *population {
	if infected > len {
		infected = len
	}

	pats := make(map[int]*patient, len)
	for i := 0; i < len; i++ {
		pats[i] = &patient{
			maxInteractions: maxInteractions,
			infectionChance: rate,
		}
	}

	pop := &population{
		patients: pats,
	}

	// infect random patients.
	for i := 0; i < infected; i++ {
		pop.infectRandom()
	}

	return pop
}

// infectRandom recursivly infects random uninfected pacients.
func (p *population) infectRandom() {
	pos := rand.Intn(len(p.patients))
	if pat := p.patients[pos]; !pat.infected {
		pat.infected = true
		return
	} else {
		p.infectRandom()
	}
}

func (p *population) loadPopulation() map[int]*patient {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.patients
}

func (p *population) removeKey(key int) {
	p.mu.Lock()
	delete(p.patients, key)
	p.mu.Unlock()
}

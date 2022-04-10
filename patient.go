package main

import (
	"math"
	"math/rand"
	"sync"
)

type patient struct {
	// maxInteractions is the max number of interactions the patient can have.
	maxInteractions int

	// interactionMu protects all the fields under.
	interactionMu sync.Mutex
	// interactedWith holds all the patients with which the patient has had contact
	// with.
	interactedWith []*patient
	// infected indicates if the patient is infected.
	infected bool
	// infectionChance is the chance that the patient gets infected.
	// should be hold as a percentage ie: 0.35, 1, 0.7
	infectionChance float32
}

func (p *patient) interact(with *patient) {
	if !p.canInteract(with) {
		return
	}

	p.interactionMu.Lock()
	p.interactedWith = append(p.interactedWith, with)
	p.interactionMu.Unlock()

	if with.isInfected() {
		if p.isInfected() {
			with.interact(p)
			return
		}

		p.infect()
	}

	// call interaction with next patient, interactions go both ways.
	with.interact(p)
}

// isInfected is a conccurency safe way of accessing p.interactedWith.
func (p *patient) canInteract(with *patient) bool {
	p.interactionMu.Lock()
	defer p.interactionMu.Unlock()

	// patient cant interact no more.
	if len(p.interactedWith) > p.maxInteractions {
		return false
	}

	// search for (with) linearly.
	for _, pat := range p.interactedWith {
		if pat == with {
			return false // alreay interacted.
		}
	}

	return true
}

// isInfected is a conccurency safe way of accessing p.infected.
func (p *patient) isInfected() bool {
	p.interactionMu.Lock()
	defer p.interactionMu.Unlock()
	return p.infected
}

func (p *patient) infect() {
	prob := math.Floor(rand.Float64()*100) / 100
	if float32(prob) < p.infectionChance {
		p.interactionMu.Lock()
		p.infected = true
		p.interactionMu.Unlock()
	}
}

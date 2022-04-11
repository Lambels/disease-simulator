package main

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

type pacient struct {
	// maxInteractions is the max number of interactions the patient can have.
	maxInteractions int

	// interactionMu protects all the fields under.
	interactionMu sync.Mutex
	// interactedWith holds all the patients with which the patient has had contact
	// with.
	interactedWith []*pacient
	// infected indicates if the patient is infected.
	infected bool
	// infectionChance is the chance that the patient gets infected.
	// should be hold as a percentage ie: 0.35, 1, 0.7
	infectionChance float32
}

func (p *pacient) interact(with *pacient) error {
	if !p.canInteract(with) {
		return errors.New("pacient cant interact with provided pacient")
	}

	p.interactionMu.Lock()
	p.interactedWith = append(p.interactedWith, with)
	p.interactionMu.Unlock()

	if with.isInfected() {
		if p.isInfected() {
			return nil
		}

		p.infect()
	}
	return nil
}

// isInfected is a conccurency safe way of accessing p.interactedWith.
func (p *pacient) canInteract(with *pacient) bool {
	p.interactionMu.Lock()
	defer p.interactionMu.Unlock()

	// cant interact with self.
	if p == with {
		return false
	}

	// patient cant interact no more.
	if len(p.interactedWith) > p.maxInteractions-1 {
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
func (p *pacient) isInfected() bool {
	p.interactionMu.Lock()
	defer p.interactionMu.Unlock()
	return p.infected
}

func (p *pacient) infect() {
	prob := math.Floor(rand.Float64()*100) / 100
	if float32(prob) < p.infectionChance {
		p.interactionMu.Lock()
		p.infected = true
		p.interactionMu.Unlock()
	}
}

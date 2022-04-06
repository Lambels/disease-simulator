package main

import "sync"

type patient struct {
	// maxInteractions is the max number of interactions the patient can have.
	maxInteractions int

	// interactionMu protects all the fields under.
	interactionMu sync.RWMutex
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
	if !p.canInteract() {
		return
	}

	p.interactionMu.Lock()
	p.interactedWith = append(p.interactedWith, with)
	p.interactionMu.Unlock()

	if with.isInfected() {

	}
}

// isInfected is a conccurency safe way of accessing p.interactedWith.
func (p *patient) canInteract() bool {
	p.interactionMu.Lock()
	defer p.interactionMu.Unlock()
	return len(p.interactedWith) < p.maxInteractions
}

// isInfected is a conccurency safe way of accessing p.infected.
func (p *patient) isInfected() bool {
	p.interactionMu.Lock()
	defer p.interactionMu.Unlock()
	return p.infected
}

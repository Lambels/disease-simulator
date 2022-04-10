package main

import (
	"log"
	"math/rand"
	"sync"
)

type population struct {
	mu sync.RWMutex

	// both clean and dirty stores store the same pointers to patients initially,
	// but the clean store doesent change its content.
	cleanStore map[int]*patient
	dirtyStore map[int]*patient

	// keys is a slice of the keys in the dirty store used to provide the random behaviour.
	keys []int

	// keyIndexMap used when removing an item from dirtyStore we can also quickly remove
	// it from keys slice above providing an O(1) complexity.
	keyIndexMap map[int]int
}

func NewPopulation(len, infected, maxInteractions int, rate float32) *population {
	if infected > len {
		infected = len
	}

	// init population.
	clean := make(map[int]*patient, len)
	dirty := make(map[int]*patient, len)
	keys := make([]int, len)
	keyIndexMap := make(map[int]int, len)
	for i := 0; i < len; i++ {
		pat := &patient{
			maxInteractions: maxInteractions,
			infectionChance: rate,
		}
		clean[i] = pat
		dirty[i] = pat
		keys[i] = i
		keyIndexMap[i] = i
	}

	pop := &population{
		cleanStore:  clean,
		dirtyStore:  dirty,
		keys:        keys,
		keyIndexMap: keyIndexMap,
	}

	// infect random patients.
	for i := 0; i < infected; i++ {
		pop.infectRandom()
	}

	return pop
}

// infectRandom recursivly infects random uninfected pacients.
//
// to only be called in initializor function.
func (p *population) infectRandom() {
	randIndex := rand.Intn(len(p.keys))

	if pat := p.dirtyStore[p.keys[randIndex]]; !pat.infected {
		pat.infected = true
		return
	} else {
		p.infectRandom()
	}
}

// RandomTwo returns 2 random pacients from the population.
func (p *population) Random() (*patient, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.dirtyStore) < 2 {
		return nil, -1
	}

	randInd1 := rand.Intn(len(p.keys))
	key := p.keys[randInd1]
	return p.dirtyStore[key], key
}

// RemoveKey removes patient with key.
//
// no-op if not found.
func (p *population) RemoveKey(key int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	index, ok := p.keyIndexMap[key]
	if !ok {
		// key not found.
		return
	}
	delete(p.keyIndexMap, key)

	wasLastIndex := len(p.keys)-1 == index

	p.keys[index] = p.keys[len(p.keys)-1]
	p.keys = p.keys[:len(p.keys)-1]

	if !wasLastIndex {
		otherKey := p.keys[index]
		p.keyIndexMap[otherKey] = index
	}

	delete(p.dirtyStore, key)
}

func (p *population) Print() {
	p.mu.RLock()
	for k, v := range p.cleanStore {
		log.Println("Pacient:", k, "Infected:", v.isInfected())
	}
	p.mu.RUnlock()
}

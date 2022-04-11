package main

import (
	"math/rand"
	"sync"

	"github.com/rodaine/table"
)

type population struct {
	mu sync.RWMutex

	// both clean and dirty stores store the same pointers to patients initially,
	// but the clean store doesent change its content.
	cleanStore map[int]*pacient
	dirtyStore map[int]*pacient

	// keys is a slice of the keys in the dirty store used to provide the random behaviour.
	keys []int

	// keyIndexMap used when removing an item from dirtyStore we can also quickly remove
	// it from keys slice above providing an O(1) complexity.
	keyIndexMap map[int]int

	maxInteractions int
}

func NewPopulation(len, infected, maxInteractions int, rate float32) *population {
	if infected > len {
		infected = len
	}

	// init population.
	clean := make(map[int]*pacient, len)
	dirty := make(map[int]*pacient, len)
	keys := make([]int, len)
	keyIndexMap := make(map[int]int, len)
	for i := 0; i < len; i++ {
		pat := &pacient{
			maxInteractions: maxInteractions,
			infectionChance: rate,
		}
		clean[i] = pat
		dirty[i] = pat
		keys[i] = i
		keyIndexMap[i] = i
	}

	pop := &population{
		cleanStore:      clean,
		dirtyStore:      dirty,
		keys:            keys,
		keyIndexMap:     keyIndexMap,
		maxInteractions: maxInteractions,
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
func (p *population) Random() (*pacient, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.dirtyStore) <= p.maxInteractions {
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

func (p *population) Analytics(infected int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	lenClean := len(p.cleanStore)
	lenDirty := len(p.dirtyStore)

	tbl := table.New("Total", "Uninfected", "Infected", "Procentage")
	widgets := []interface{}{}

	widgets = append(widgets, lenClean)
	if lenDirty == lenClean {
		// before simulation.
		widgets = append(widgets, lenClean-infected)
		widgets = append(widgets, infected)
		widgets = append(widgets, (infected/lenClean)*100)
	} else {
		// after simulation.
		var infectionCount int
		for _, pat := range p.cleanStore {
			if pat.isInfected() {
				infectionCount++
			}
		}

		widgets = append(widgets, lenClean-infectionCount)
		widgets = append(widgets, infectionCount)
		widgets = append(widgets, (infectionCount/lenClean)*100)
	}
	tbl.AddRow(widgets...)
	tbl.Print()
}

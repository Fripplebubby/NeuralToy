package main

import (
		"math/rand"
		)


type Genome struct {
	weights []float64
	weightCount int
	fitness float64
}

func NewGenome(weightCount int) Genome {
	gen := Genome{make([]float64, weightCount), weightCount, -1000}
	for i := range gen.weights {
		gen.weights[i] = rand.Float64()
	}
	return gen
}

func (gen *Genome) Copy() Genome {
	gen2 := Genome{
		make([]float64, len(gen.weights)),
		gen.weightCount,
		gen.fitness,
	}
	copy(gen2.weights, gen.weights)
	return gen2
}

type Population struct {
	members []Genome
	size int
}

func NewPopulation() Population {
	return Population {
		make([]Genome, 0),
		0,
	}
}

func (pop *Population) Add(g Genome){
	pop.members = append(pop.members, g)
	pop.size++
}

func (pop *Population) Set(g []Genome) {
	pop.members = make([]Genome, len(g))
	copy(pop.members, g)
	pop.size = len(g)
}

func (pop *Population) Get(i int) *Genome {
	return &pop.members[i]
}

func (pop *Population) GetRange(i, j int) []Genome {
	return pop.members[i:j]
}

func (pop Population) Len() int {
	return pop.size
}

func (pop Population) Less(i, j int) bool {
	return pop.members[i].fitness < pop.members[j].fitness
} 

func (pop Population) Swap(i, j int) {
	g2 := pop.members[i].Copy()
	pop.members[i] = pop.members[j].Copy()
	pop.members[j] = g2
}


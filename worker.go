package main

import (
	"math"
	"sort"
	"math/rand"
	"fmt"
)

type Worker struct {
	id int
	nextRound chan Round 
	newWeights chan Genome
	newRounds chan []Round
	pop Population
	rounds []Round
	fittestHappiestMostProductive Genome
	done chan bool
}

func NewWorker(id int) Worker {
	pop := NewPopulation()
	for i := 0; i < PopSize; i++ {
		pop.Add(NewGenome(WeightCount))
	}

	return Worker {
		id, 
		make (chan Round),
		make (chan Genome),
		make (chan []Round),
		pop,
		make ([]Round, 0),
		NewGenome(0),
		make (chan bool),
	}
}

func (w *Worker) Work(complete chan<- *Worker){
	for {
		select {
			case rnd := <-w.nextRound:
				w.ProduceGenome(rnd)
				complete <- w
			case weights := <-w.newWeights:
				w.NewBest(weights)
			case <-w.done:
				complete <- w
			case rnds := <-w.newRounds:
				w.InjectRounds(rnds)
		}
	}
}

func (w *Worker) NewBest(bestGen Genome){
	for _, p := range w.pop.members {
		p.weights = bestGen.weights
		// might be necessary to copy fitness, might not
		p.fitness = bestGen.fitness
	}
}

func (w *Worker) InjectRounds(rnd []Round){
	w.rounds = make([]Round, len(rnd))
	for i := range w.rounds {
		w.rounds[i] = rnd[i].Copy()
	}
}

func (w *Worker) ProduceGenome(r Round){
	w.rounds = append(w.rounds, r)
	in, out := buildInput(w.rounds)
	net := NewNeuralNet(NumInputs, NumOutputs, NumHiddenLayers, NumNeuronsPerHiddenLayer)
	k := 0
	for bestFit := -1000.0; bestFit < NumFitnessGoal; k++ {
		select {
			case <-w.done :
				bestFit = 1
				return 
			
			default:
				w.calcFitness(in, out, net)
				w.pruneWeaklings()
				bestFit = w.fittestHappiestMostProductive.fitness
				if k % 1000 == 0 {
					fmt.Printf("Worker %v new best: %v\tSigmoid: %v\t%v + %v = %v\n", w.id + 1, bestFit, Sigmoid(bestFit), r.Int1, r.Int2, w.TestGenome(r))
				}

		}
		
	}
}

func (w *Worker) GetBestGenome() Genome {
	return w.fittestHappiestMostProductive
}

func (w *Worker) DropWork(){
	w.done <- true
}

func buildInput(rnds []Round) (in, out [][]float64) {
	input := make([][]float64, 0)
	expectedOutput := make([][]float64, 0)
	for i := len(rnds) - 1; i >= 0; i -- {
		input = append(input,  rnds[i].Input)
		expectedOutput = append(expectedOutput, rnds[i].ExpectedOutput)
	}
	return input, expectedOutput
}

func (w *Worker) calcFitness(input, expectedOutput [][]float64, net NeuralNet) {
	for i := range w.pop.members {
		genome := w.pop.Get(i)
		genome.fitness = 0
		net.PutWeights(genome.weights)
		for j := range input {
			output := net.Update(input[j]) 
			for k := range output {
				genome.fitness += -math.Abs(output[k]-expectedOutput[j][k])
			}
		}
		genome.fitness = Sigmoid(genome.fitness)
	}
}

func (w *Worker) pruneWeaklings(){

	sort.Sort(sort.Reverse(w.pop))
	keepers := w.pop.members[0:PopCutoff]
	w.pop.Set(keepers)

	keepersInd := 0
	for i := PopCutoff; i < PopSize; i++ {
		oldGen := keepers[keepersInd % PopCutoff]
		newgen := oldGen.Copy()
		newgen.Mutate()
		w.pop.Add(newgen)
		keepersInd++
	}
	w.fittestHappiestMostProductive = w.pop.members[0]
}

func (w *Worker) TestGenome(rnd Round) int {
	net := NewNeuralNet(NumInputs, NumOutputs, NumHiddenLayers, NumNeuronsPerHiddenLayer)
	net.PutWeights(w.fittestHappiestMostProductive.weights)
	output := net.Update(rnd.Input)
	answer := 0
	for i := range output {
		if math.Floor(output[i] + 0.5) >= 1 {
			answer += i
		}
	}
	return answer
}

func (gen *Genome) Mutate() {
	for i := range gen.weights {
		if rand.Float64() < MutationRate {
			chng := rand.Float64()
			neg := rand.Float64()
			if neg > 0.5 {
				chng *= -1
			}
			gen.weights[i] += chng
		}
	} 
}

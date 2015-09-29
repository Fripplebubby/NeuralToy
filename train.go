package main

import (
	"math/rand"
	"fmt"
	"encoding/json"
    	"io/ioutil"
)

type Trainer struct {
	workers []Worker
	nextWorker chan *Worker
	trRounds int
	teRounds int
}

func NewTrainer(trRounds, teRounds, weightCount int) Trainer {
	workers := make([]Worker, CPUCores)
	for i := range workers {
		workers[i] = NewWorker(i, weightCount)
	}

	nextWorker := make(chan *Worker, CPUCores)	

	for i := range workers {
		go workers[i].Work(nextWorker)
	}

	return Trainer {
		workers,
		nextWorker,
		trRounds,
		teRounds,
	}
}

func (train *Trainer) LoadGenome(){
	jsonGen, err := ioutil.ReadFile(StrPath + GenomeReadName)
	if err != nil {
		panic(err)
	}

	var weights []float64
	jsonErr := json.Unmarshal(jsonGen, &weights)

	if jsonErr != nil {
		panic(jsonErr)
	}

	gen := NewGenome(len(weights))
	gen.weights = weights

	for _, w := range train.workers {
		w.NewBest(gen.Copy())
	}

	jsonRounds, err2 := ioutil.ReadFile(StrPath + RoundsReadName)
	if err2 != nil {
		panic(err2)
	}

	var rounds []Round
	jsonErr2 := json.Unmarshal(jsonRounds, &rounds)

	if jsonErr2 != nil {
		panic(jsonErr2)
	}

	for _, w := range train.workers {
		w.newRounds <- rounds
	}
}

func (train *Trainer) SaveGenome(){
	write(train.workers[0].fittestHappiestMostProductive.weights, GenomeWriteName)
	write(train.workers[0].rounds, RoundsWriteName)
	fmt.Printf("Saved round %v\n", len(train.workers[0].rounds))
}

func (train *Trainer) Generate(){
	for onRound := 0; onRound < train.trRounds; onRound++ {
		nextRound := NewRound()
		fmt.Println("Dispensing round...")
		for i := range train.workers {
			train.workers[i].nextRound <- nextRound 
		}

		w := <- train.nextWorker 
		skipId := w.id
		fit := w.GetBestGenome()
		// copy the best genome out to every worker 
		for _, wrk := range train.workers {
			select {
				case <- train.nextWorker:
					wrk.newWeights <- fit.Copy()
					continue
				default: 
					if wrk.id != skipId {
						wrk.DropWork() // does not terminate the work loop, but will reset all workers to waiting for the next round
						<- train.nextWorker // clear the done queue
						wrk.newWeights <- fit.Copy()
					}
			}			
		}
		fmt.Println("Received trained response!")
		train.SaveGenome()
	}

	fmt.Println("BEGIN TESTING")
	for test := 0; test < train.teRounds; test++ {
		nextRound := NewRound()
		result := train.workers[0].TestGenome(nextRound) 
		fmt.Printf("%v + %v = %v\n", nextRound.Int1, nextRound.Int2, result)
	}

}

type Round struct {		
	Int1 int
	Int2 int
	Answer int
	Input []float64 // 				[0,0,1,0,1,0,0,0,0,0]  			2 + 4 
	ExpectedOutput []float64 // 	[0,0,0,0,0,0,1,0,0,0]			= 6 
}

func NewRound() Round {
	Int1 := rand.Intn(6)
	Int2 := rand.Intn(5)
	input := intsToInputArray(Int1, Int2)
	input[10] = float64(rand.Intn(2)) // multiplication or addition?
	var output []float64
	if input[10] == 1 {
		output = intToInputArray(Int1 * Int2)
	} else {
		output = intToInputArray(Int1 + Int2)
	}
	return Round {
		Int1,
		Int2,
		Int1 + Int2,
		input,
		output,
	}
}

func (r *Round) Copy() Round {
	rnd := NewRound()
	rnd.Int1 = r.Int1
	rnd.Int2 = r.Int2
	rnd.Answer = r.Answer
	copy(rnd.Input, r.Input)
	copy(rnd.ExpectedOutput, r.ExpectedOutput)
	return rnd
}

// [0,0,0,0,1,0,0,0,0,0] = 4 
func intsToInputArray(in, in2 int) []float64{
	toRet := make([]float64, NumInputs)
	toRet[in] += 1
	toRet[in2] += 1
	return toRet
}

func intToInputArray(in int) []float64{
	toRet := make([]float64, NumOutputs)
	toRet[in] += 1
	return toRet
}

func write(data interface{}, name string){
	jsonGen, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	writeErr := ioutil.WriteFile(StrPath + name, jsonGen, 0666)
	if writeErr != nil {
		panic(writeErr)
	}
}

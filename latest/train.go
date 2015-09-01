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

func NewTrainer(trRounds, teRounds int) Trainer {
	workers := make([]Worker, CPUCores)
	for i := range workers {
		workers[i] = NewWorker(i)
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

func (train *Trainer) LoadGenome(path string){
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var weights []float64
	jsonErr := json.Unmarshal(f, &weights)

	if jsonErr != nil {
		panic(jsonErr)
	}

	gen := NewGenome(len(weights))
	gen.weights = weights

	for _, w := range train.workers {
		w.NewBest(gen.Copy())
	}


}

func (train *Trainer) SaveGenome(path string){


	jsonGen, err := json.Marshal(train.workers[0].GetBestGenome().weights)

	if err != nil {
		panic(err)
	}

	writeErr := ioutil.WriteFile(path, jsonGen, 0666)
	if writeErr != nil {
		panic(writeErr)
	}
	
	testData := train.workers[0].rounds[0].Serialize()

	/*jTest, jtErr := json.Marshal()
	if jtErr != nil {
		panic(jtErr)
	}*/
	fmt.Printf("%v\n",testData)
}

func (train *Trainer) Generate(){
	for onRound := 0; onRound < train.trRounds; onRound++ {
		nextRound := NewRound()
		fmt.Println("Dispensing round...")
		for i := range train.workers {
			//fmt.Printf("Sending round to #%v\n",i)
			train.workers[i].nextRound <- nextRound 
		}

		w := <- train.nextWorker 
		//fmt.Printf("%v is in first!\n", w.id)
		skipId := w.id
		fit := w.GetBestGenome()
		// copy the best genome out to every worker 
		for _, wrk := range train.workers {
			select {
				case <- train.nextWorker:
					//fmt.Printf("Heard from %v\n", wrk.id)
					wrk.newWeights <- fit.Copy()
					continue
				default: 
					if wrk.id != skipId {
						//fmt.Printf("Sending stop to %v\n", wrk.id)
						wrk.DropWork() // does not terminate the work loop, but will reset all workers to waiting for the next round
						<- train.nextWorker // clear the done queue
						//fmt.Printf("Heard from %v\n", wrk.id)
						wrk.newWeights <- fit.Copy()
					}
			}			
		}
		fmt.Println("Received trained response!")
	}

	fmt.Println("BEGIN TESTING")
	for test := 0; test < train.teRounds; test++ {
		nextRound := NewRound()
		result := train.workers[0].TestGenome(nextRound) 
		fmt.Printf("%v + %v = %v\n", nextRound.int1, nextRound.int2, result)
	}

}

type Round struct {		
	int1 int
	int2 int
	answer int
	input []float64 // [0,0,0,0,1,0,0,0,0,0] = 4 
	expectedOutput []float64 // same format
}

func NewRound() Round {
	int1 := rand.Intn(6)
	int2 := rand.Intn(5)
	return Round {
		int1,
		int2,
		int1 + int2,
		intsToInputArray(int1, int2),
		intToInputArray(int1 + int2),
	}
}

type SerialRound struct {
	Int1 int 
	Int2 int 
	Answer int
}

func (rnd *Round) Serialize() SerialRound {
	return SerialRound {
		rnd.int1,
		rnd.int2,
		rnd.answer,
	}
}

func (sr *SerialRound) DeSerialize() Round {
	return Round {
		sr.Int1,
		sr.Int2,
		sr.Answer,
		intsToInputArray(sr.Int1, sr.Int2),
		intToInputArray(sr.Answer),
	}
}

// [0,0,0,0,1,0,0,0,0,0] = 4 
func intsToInputArray(in, in2 int) []float64{
	toRet := []float64{0,0,0,0,0,0,0,0,0,0}
	toRet[in] += 1
	toRet[in2] += 1
	return toRet
}

func intToInputArray(in int) []float64{
	toRet := []float64{0,0,0,0,0,0,0,0,0,0}
	toRet[in] += 1
	return toRet
}

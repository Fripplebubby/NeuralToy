package main

import (
	"testing"
	"math"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"math/rand"
)

func TestGenomes(*testing.T){
	rand.Seed(time.Now().UTC().UnixNano())
	toTest := LoadGenome()
	net := NewNeuralNet(NumInputs, NumOutputs, NumHiddenLayers, NumNeuronsPerHiddenLayer)
	for i := 0; i < NumTestingRounds; i++ {
		nextRound := NewRound()
		net.PutWeights(toTest.weights)
		result := net.Update(nextRound.Input)
		answer := 0
		for i, r := range result {
			if math.Floor(r + 0.5) == 1 {
				answer += i
			}
		}
		fmt.Printf("%v + %v = %v\n", nextRound.Int1, nextRound.Int2, answer)
	}
}

func LoadGenome() Genome {
	f, err := ioutil.ReadFile(StrPath + GenomeReadName)
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
	return gen
}

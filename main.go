package main

import (
		"math/rand"
		"time"
		"runtime"
		//"github.com/davecheney/profile"
		//"fmt"
		)



func main() {
	runtime.GOMAXPROCS(CPUCores)
	//defer profile.Start(profile.CPUProfile).Stop()
	

	rand.Seed(time.Now().UTC().UnixNano())
	
	
	net := NewNeuralNet(NumInputs, NumOutputs, NumHiddenLayers, NumNeuronsPerHiddenLayer)
	WeightCount = net.GetNumberOfWeights()

	trainer := NewTrainer(NumTrainingRounds, NumTestingRounds)
	trainer.Generate()
}	

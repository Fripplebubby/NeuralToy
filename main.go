package main

import (
		"math/rand"
		"time"
		"runtime"
		//"github.com/davecheney/profile"
		//"fmt"
		)



func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//defer profile.Start(profile.CPUProfile).Stop()
	
	rand.Seed(time.Now().UTC().UnixNano())
	
	net := NewNeuralNet(NumInputs, NumOutputs, NumHiddenLayers, NumNeuronsPerHiddenLayer)

	trainer := NewTrainer(NumTrainingRounds, NumTestingRounds, len(net.weights))
	trainer.LoadGenome()
	trainer.Generate()
	trainer.SaveGenome()
}	

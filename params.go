package main

const PopSize = 150
const PopCutoff = 25
const MutationRate = 0.2
const NumInputs = 10
const NumHiddenLayers = 10
const NumNeuronsPerHiddenLayer = 15
const NumOutputs = 10
const NeuronBias = -1
const NumTrainingRounds = 1
const NumTestingRounds = 0

var NumFitnessGoal = 0.5

var WeightCount int // set at runtime 
var CPUCores int 
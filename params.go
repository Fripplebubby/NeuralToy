package main

const CPUCores = 6
const PopSize = 150
const PopCutoff = 25
const MutationRate = 0.2
const NumInputs = 3
const NumHiddenLayers = 3
const NumNeuronsPerHiddenLayer = 2
const NumOutputs = 3
const NeuronBias = -1
const NumTrainingRounds = 15
const NumTestingRounds = 10

var NumFitnessGoal = 0.5

var WeightCount int // set at runtime 
package main

const PopSize = 150
const PopCutoff = 25
const MutationRate = 0.2
const NumInputs = 10
const NumHiddenLayers = 10
const NumNeuronsPerHiddenLayer = 15
const NumOutputs = 10
const NeuronBias = -1
const NumTrainingRounds = 15
const NumTestingRounds = 15

const StrPath = "/home/ec2-user/NeuralToy/"
const GenomeWriteName = "Genome1"
const RoundsWriteName = "Rounds1"
const GenomeReadName = "Genome1"
const RoundsReadName = "Rounds1"

var NumFitnessGoal = 0.001

var WeightCount int // set at runtime 
var CPUCores int 

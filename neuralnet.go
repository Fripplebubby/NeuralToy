package main

import (
	"math/rand"
	"math"
)

type Neuron struct {
	numInputs int
	weights []float64
}

func NewNeuron(inputs int) Neuron {
	neur := Neuron{inputs, make([]float64, inputs)}
	for i := range neur.weights {
		neur.weights[i] = rand.Float64()
	}
	return neur
}

type NeuronLayer struct {
	numNeurons int
	neurons []Neuron
}

func NewNeuronLayer(numNeurons int, numInputsPerNeuron int) NeuronLayer {
	layer := NeuronLayer{numNeurons, make([]Neuron, numNeurons)}
	for i := range layer.neurons {
		layer.neurons[i] = NewNeuron(numInputsPerNeuron)
	}
	return layer
}

type NeuralNet struct {
	numInputs int
	numOutputs int
	numHiddenLayers int 
	numNeuronsPerLayer int
	layers []NeuronLayer
	weights []float64
}

func NewNeuralNet(numInputs int, numOutputs int, numHiddenLayers int, numNeuronsPerLayer int) NeuralNet {
	net := NeuralNet{numInputs,numOutputs,numHiddenLayers,numNeuronsPerLayer, make([]NeuronLayer, numHiddenLayers + 2), make([]float64, 0)}
	for i := range net.layers {
		if i < numHiddenLayers + 1 {
			if i > 0 {
				// hidden layer
				net.layers[i] = NewNeuronLayer(numNeuronsPerLayer, net.layers[i-1].numNeurons + 1) // + 1 is bias
				} else {
					// input layer
					net.layers[i] = NewNeuronLayer(numInputs, 0)
				}
		} else {
			// output layer
			net.layers[i] = NewNeuronLayer(numOutputs, net.layers[i-1].numNeurons + 1) // bias again
		}
	}
	net.weights = make([]float64, net.GetNumberOfWeights())
	return net
}

func (nn NeuralNet) GetWeights() []float64 {
	/*length := nn.numInputs + nn.numOutputs + (nn.numHiddenLayers * nn.numNeuronsPerLayer)
	ret := make([]float64, length)
	i := 0 
	for j := range nn.layers {
		for k := range nn.layers[j].neurons {
			for l := range nn.layers[j].neurons[k].weights {
				ret[i] = nn.layers[j].neurons[k].weights[l]
				i++ 
			}
		}
	}
	return ret*/
	return nn.weights
}

func (nn NeuralNet) GetNumberOfWeights() int {
	i := 0
	for j := range nn.layers {
		for k := range nn.layers[j].neurons {
			for range nn.layers[j].neurons[k].weights {
				i++
			}
		}
	}
	return i
}

func (nn *NeuralNet) PutWeights(weights []float64) {
	/*i := 0
	for j := range nn.layers {
		for k := range nn.layers[j].neurons {
			for l := range nn.layers[j].neurons[k].weights {
				nn.layers[j].neurons[k].weights[l] = weights[i]
				i++
			}
		}
	}*/
	copy(nn.weights, weights)
	//nn.weights = weights
}

func (nn NeuralNet) Update(inputs []float64) []float64 {
	var cWeight int 
	outputs := make([]float64, 0)
	for i := 1; i < nn.numHiddenLayers + 2; i++ {
		if i > 1 {
			inputs = outputs
		}
		outputs = make([]float64, nn.layers[i].numNeurons)
		cWeight = 0
		for j := range nn.layers[i].neurons {

			netInput := 0.0

			for k := 0; k < nn.layers[i].neurons[j].numInputs-1; k++ {
				index := ((nn.numHiddenLayers + 2) * i) + (nn.numNeuronsPerLayer * j) + k

				netInput += nn.weights[index] * inputs[cWeight]
				cWeight++
			}

			// bias
			if nn.layers[i].neurons[j].numInputs > 0 {
				index := ((nn.numHiddenLayers + 2) * i) + (nn.numNeuronsPerLayer * j) + (nn.numInputs-1)
				netInput += nn.weights[index] * NeuronBias
			} 


			outputs[j] = Sigmoid(netInput)
			cWeight = 0

		}
	}
	return outputs
}

func Sigmoid(input float64) float64 {
	return 1 / (1 + math.Pow(math.E, -input))
}
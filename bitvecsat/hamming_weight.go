package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
	"fmt"
)

type HammingWeightConstrain struct {
	AIndex int
	WeightIndex int

	aResultIndex int
	bResultIndex int
}

func (constrain HammingWeightConstrain) Materialize(problem *Problem) []sat.Clause {
	return VectorsAreEqual(problem.Vectors[constrain.aResultIndex], problem.Vectors[constrain.bResultIndex])
}

func (constrain HammingWeightConstrain) AddToProblem(problem *Problem) {
	// (1 << width of WeightIndex) - 1 is the maximum amount of bits
	// representable by WeightIndex.
	a := problem.Vectors[constrain.AIndex]
	weight := problem.Vectors[constrain.WeightIndex]

	maxBitsInA := uint((1 << weight.Width) - 1)
	if a.Width > maxBitsInA {
		panic(fmt.Sprintf("weight vector (%d) not large enough for possible # of bits (%d)",
			weight.Width, a.Width))
	}

	widthOfResults := uint(1)
	resultIndices := make([]int, a.Width)
	for i := uint(0); i < a.Width; i++ {
		index := a.SatVarIndices[i]
		resultIndices[i] = problem.AddBoundVector(1, []int{index})
	}

	for len(resultIndices) > 1 {
		nextResultIndices := make([]int, (len(resultIndices) / 2) + (len(resultIndices) % 2))

		for i := 0; i < len(resultIndices) / 2; i++ {
			nextResultIndices[i] = problem.AddNewVector(widthOfResults + 1)

			a := resultIndices[i * 2]
			b := resultIndices[i * 2 + 1]

			extendedA := problem.AddLeftExtendedVector(a, widthOfResults + 1)
			extendedB := problem.AddLeftExtendedVector(b, widthOfResults + 1)
			PlusConstrain{
				AIndex: extendedA,
				BIndex: extendedB,
				SumIndex: nextResultIndices[i],
			}.AddToProblem(problem)
		}
		if len(resultIndices) % 2 == 1 {
			nextResultIndices[len(resultIndices) / 2] = problem.AddLeftExtendedVector(resultIndices[len(resultIndices) - 1], widthOfResults + 1)
		}
		resultIndices = nextResultIndices
		widthOfResults++
	}

	// TODO: hack
	finalResult := problem.Vectors[resultIndices[0]]
	extWidth := weight.Width
	if extWidth < finalResult.Width {
		extWidth = finalResult.Width
	}
	constrain.aResultIndex = problem.AddLeftExtendedVector(resultIndices[0], extWidth)
	constrain.bResultIndex = problem.AddLeftExtendedVector(constrain.WeightIndex, extWidth)
	problem.AddNewConstrain(constrain)
}

func (constrain HammingWeightConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	return "hamming_weight"
}

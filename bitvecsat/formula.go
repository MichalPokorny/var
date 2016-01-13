package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
)

// Term can:
//   * return a vector,
//   * partially materialize constraints

// Constraints:
//   v = sum(v1, v2, ...)
//   v = product(v1, v2, ...)

type Vector struct {
	// May be nil if the vector's not materialized yet.
	// TODO: variable width (currently 8)
	SatVarIndices []int // [0] is least significant bit
	Width uint
}

func (problem *Problem) GetBitsInAssignment(assignment sat.Assignment, vectorIndex int) string {
	v := ""
	for i := len(problem.Vectors[vectorIndex].SatVarIndices) - 1; i >= 0; i-- {
		bit := assignment[problem.Vectors[vectorIndex].SatVarIndices[i]]
		if bit {
			v += "1"
		} else {
			v += "0"
		}
	}
	return v
}

func (problem *Problem) GetValueInAssignment(assignment sat.Assignment, vectorIndex int) int {
	value := 0
	vector := problem.Vectors[vectorIndex]
	for i := len(vector.SatVarIndices) - 1; i >= 0; i-- {
		value = value << 1
		if assignment[vector.SatVarIndices[i]] {
			value = value | 1
		}
	}
	return value
}

// func Sum(inputs, output) Constrain
// func And(inputs, output) Constrain

type Constrain interface {
	Materialize(problem *Problem) []sat.Clause
	Dump(problem *Problem, assignment sat.Assignment) string
}

type Problem struct {
	Vectors []Vector
	Constrains []Constrain

	LastSatVarIndex int
	AlwaysZeroBitIndex *int
}

func (problem *Problem) AddNewVector(width uint) int {
	vector := Vector{
		Width: width,
		SatVarIndices: make([]int, width),
	}
	for i := 0; i < int(width); i++ {
		vector.SatVarIndices[i] = problem.LastSatVarIndex
		problem.LastSatVarIndex++
	}
	problem.Vectors = append(problem.Vectors, vector)
	return len(problem.Vectors) - 1
}

func (problem *Problem) AddBoundVector(width uint, indices []int) int {
	vector := Vector{
		Width: width,
		SatVarIndices: indices,
	}
	if width != uint(len(indices)) {
		panic("bad width")
	}
	problem.Vectors = append(problem.Vectors, vector)
	return len(problem.Vectors) - 1
}

func (problem *Problem) AddConstantVector(width uint, value int) int {
	index := problem.AddNewVector(width)
	problem.AddNewConstrain(LiteralConstrain{AIndex: index, Value: value})
	return index
}

func (problem *Problem) AddNewConstrain(constrain Constrain) {
	problem.Constrains = append(problem.Constrains, constrain)
}

func (problem *Problem) MakeSatFormula() sat.Formula {
	clauses := make([]sat.Clause, 0)
	for i := 0; i < len(problem.Constrains); i++ {
		clauses = append(clauses, problem.Constrains[i].Materialize(problem)...)
	}
	return sat.Formula{Clauses: clauses}
}

func (problem *Problem) GetAlwaysZeroBitIndex() int {
	if problem.AlwaysZeroBitIndex != nil {
		return *problem.AlwaysZeroBitIndex
	}

	alwaysZeroVectorIndex := problem.AddConstantVector(1, 0)
	problem.AlwaysZeroBitIndex = new(int)
	*problem.AlwaysZeroBitIndex = problem.Vectors[alwaysZeroVectorIndex].SatVarIndices[0]
	return *problem.AlwaysZeroBitIndex
}

func (problem *Problem) AddLeftExtendedVector(index int, width uint) int {
	original := problem.Vectors[index]
	if original.Width > width {
		panic("left extension is not an extension")
	}
	vars := original.SatVarIndices
	for i := uint(0); i < (width - original.Width); i++ {
		vars = append(vars, problem.GetAlwaysZeroBitIndex())
	}
	return problem.AddBoundVector(width, vars)
}

// width of each vector

//type Term struct {
//	// SUM(Term1, Term2, ...)
//	// PRODUCT(Term1, Term2, ...)
//	// OR(Term1, Term2, ...)
//	// AND(Term1, Term2, ...)
//	// XOR(Term1, Term2, ...)
//	// NOT(Term1)
//}
//
//type Formula struct {
//	// AND(Formula, Formula, ...)
//	// OR(Formula, Formula, ...)
//	// RelOp(Term, Term): < <= = != => >
//}

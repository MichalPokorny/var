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
	for i := len(problem.Vectors[vectorIndex].SatVarIndices) - 1; i >= 0; i-- {
		value = value << 1
		if assignment[problem.Vectors[vectorIndex].SatVarIndices[i]] {
			value = value | 1
		}
	}
	return value
}

// func Sum(inputs, output) Constrain
// func And(inputs, output) Constrain

type Constrain interface {
	AddToProblem(problem *Problem)
	// TODO: partial materialization
	Materialize(problem *Problem) []sat.Clause
}

type Problem struct {
	Vectors []Vector
	Constrains []Constrain
}

func (problem *Problem) AddNewVector(width uint) int {
	problem.Vectors = append(problem.Vectors, Vector{Width: width})
	return len(problem.Vectors) - 1
}

func (problem *Problem) AddNewConstrain(constrain Constrain) {
	problem.Constrains = append(problem.Constrains, constrain)
}

func (problem *Problem) PrepareSat() {
	satVar := 0
	for i := 0; i < len(problem.Vectors); i++ {
		problem.Vectors[i].SatVarIndices = make([]int, problem.Vectors[i].Width)
		for j := 0; j < int(problem.Vectors[i].Width); j++ {
			problem.Vectors[i].SatVarIndices[j] = satVar
			satVar++
		}
	}
}

func (problem *Problem) MakeSatFormula() sat.Formula {
	clauses := make([]sat.Clause, 0)
	for i := 0; i < len(problem.Constrains); i++ {
		clauses = append(clauses, problem.Constrains[i].Materialize(problem)...)
	}
	return sat.Formula{Clauses: clauses}
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
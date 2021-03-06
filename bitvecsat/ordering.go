package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
	"fmt"
)

// TODO: iota
const LTE = 0;
const LT = 1;

// A <= B
type OrderingConstrain struct {
	AIndex int
	BIndex int
	Type int  // LTE or LT

	IsQuery bool  // if true, it's a query, not an assertion
	QueryResultBitIndex int

	// TODO: maybe fewer variables, bigger conditions?
	bitLtIndex int
	bitEqIndex int
	suffixMeetsIndex int
	nextBitDecidesIndex int
}

func (constrain OrderingConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	bitLt := problem.Vectors[constrain.bitLtIndex]
	bitEq := problem.Vectors[constrain.bitEqIndex]
	suffixMeets := problem.Vectors[constrain.suffixMeetsIndex]
	nextBitDecides := problem.Vectors[constrain.nextBitDecidesIndex]

	width := a.Width

	clauses := make([]sat.Clause, 0)
	clauses = append(clauses, VectorsBitwise(a, b, bitLt, sat.LtConstrain)...)
	clauses = append(clauses, VectorsBitwise(a, b, bitEq, sat.EquivConstrain)...)
	clauses = append(clauses, VectorsBitwise(bitLt, nextBitDecides, suffixMeets, sat.OrConstrain)...)

	for i := 0; i < int(width); i++ {
		// suffixLTE[i] <=> (a[0..i] <= b[0..i])
		bitEqBit := bitEq.SatVarIndices[i]
		nextBitDecidesBit := nextBitDecides.SatVarIndices[i]

		if i > 0 {
			nextSuffixLteBit := suffixMeets.SatVarIndices[i - 1]
			clauses = append(clauses, sat.AndConstrain(bitEqBit, nextSuffixLteBit, nextBitDecidesBit)...)
		} else {
			// TODO: allow easier equalities between variables -- make one variable
			if constrain.Type == LTE {
				clauses = append(clauses, sat.BitsAlwaysEqual(bitEqBit, nextBitDecidesBit)...)
			} else if constrain.Type == LT {
				clauses = append(clauses, sat.BitIsFalse(nextBitDecidesBit)...)
			} else {
				panic("unknown constrain type")
			}
		}
	}

	resultBit := suffixMeets.SatVarIndices[width - 1]
	if constrain.IsQuery {
		clauses = append(clauses, sat.BitsAlwaysEqual(resultBit, constrain.QueryResultBitIndex)...)
	} else {
		clauses = append(clauses, sat.BitIsTrue(resultBit)...)
	}
	return clauses
}

func (constrain OrderingConstrain) AddToProblem(problem *Problem) {
	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	width := a.Width
	if width != b.Width {
		panic("unequal bit widths")
	}

	constrain.bitLtIndex = problem.AddNewVector(width)
	constrain.bitEqIndex = problem.AddNewVector(width)
	constrain.suffixMeetsIndex = problem.AddNewVector(width)
	constrain.nextBitDecidesIndex = problem.AddNewVector(width)
	problem.AddNewConstrain(constrain)
}

func (constrain OrderingConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	var t string
	if constrain.Type == LTE {
		t = "<="
	} else {
		t = "<"
	}
	deciding := fmt.Sprintf("#%d[== %v] %s #%d[== %v]", constrain.AIndex, problem.GetValueInAssignment(assignment, constrain.AIndex), t, constrain.BIndex, problem.GetValueInAssignment(assignment, constrain.BIndex))
	str := "ordering("
	if constrain.IsQuery {
		str += fmt.Sprintf("[%s] <=> bit %d [== %v]", deciding, constrain.QueryResultBitIndex, assignment[constrain.QueryResultBitIndex])
	} else {
		str += deciding
	}
	str += ")"
	return str
}

package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
)

// A <= B
type LTEConstrain struct {
	AIndex int
	BIndex int

	// TODO: maybe fewer variables, bigger conditions?
	bitLtIndex int
	bitEqIndex int
	suffixLteIndex int
	nextBitDecidesIndex int
}

func (constrain LTEConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	bitLt := problem.Vectors[constrain.bitLtIndex]
	bitEq := problem.Vectors[constrain.bitEqIndex]
	suffixLte := problem.Vectors[constrain.suffixLteIndex]
	nextBitDecides := problem.Vectors[constrain.nextBitDecidesIndex]

	width := a.Width()
	if ((width != b.Width()) || (width != bitLt.Width()) || (width != bitEq.Width()) || (width != suffixLte.Width())) {
		panic("unequal bit widths")
	}

	clauses := make([]sat.Clause, 0)
	clauses = append(clauses, VectorsBitwise(a, b, bitLt, sat.LtConstrain)...)
	clauses = append(clauses, VectorsBitwise(a, b, bitEq, sat.EquivConstrain)...)
	clauses = append(clauses, VectorsBitwise(bitLt, nextBitDecides, suffixLte, sat.OrConstrain)...)

	for i := 0; i < width; i++ {
		// suffixLTE[i] <=> (a[0..i] <= b[0..i])
		bitEqBit := bitEq.SatVarIndices[i]
		nextBitDecidesBit := nextBitDecides.SatVarIndices[i]

		if i > 0 {
			nextSuffixLteBit := suffixLte.SatVarIndices[i - 1]
			clauses = append(clauses, sat.AndConstrain(bitEqBit, nextSuffixLteBit, nextBitDecidesBit)...)
		} else {
			// TODO: allow easier equalities between variables -- make one variable
			clauses = append(clauses, sat.BitsAlwaysEqual(bitEqBit, nextBitDecidesBit)...)
		}
	}

	clauses = append(clauses, sat.BitIsTrue(suffixLte.SatVarIndices[width - 1])...)
	return clauses
}

func (constrain LTEConstrain) AddToProblem(problem *Problem) {
	constrain.bitLtIndex = problem.AddNewVector()
	constrain.bitEqIndex = problem.AddNewVector()
	constrain.suffixLteIndex = problem.AddNewVector()
	constrain.nextBitDecidesIndex = problem.AddNewVector()
	problem.AddNewConstrain(constrain)
}

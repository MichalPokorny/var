package bitvecsat

import (
	"fmt"
	"github.com/MichalPokorny/var/sat"
)

type BitConstrain func(int, int, int) []sat.Clause

type BitwiseLogicalConstrain struct {
	AIndex int
	BIndex int
	YIndex int

	BitConstrain BitConstrain
}

func VectorsBitwise(a Vector, b Vector, y Vector, bitConstrain BitConstrain) []sat.Clause {
	width := a.Width
	if ((width != b.Width) || (width != y.Width)) {
		panic("unequal bit widths")
	}

	clauses := make([]sat.Clause, 0) // TODO: exact width
	for i := 0; i < int(width); i++ {
		a := a.SatVarIndices[i]
		b := b.SatVarIndices[i]
		y := y.SatVarIndices[i]
		clauses = append(clauses, bitConstrain(a, b, y)...)
	}
	return clauses
}

func VectorsAreEqual(a Vector, b Vector) []sat.Clause {
	width := a.Width
	if width != b.Width {
		panic("unequal bit widths")
	}

	// TODO: exact size
	clauses := make([]sat.Clause, 0)
	for i := 0; i < int(width); i++ {
		aBit := a.SatVarIndices[i]
		bBit := b.SatVarIndices[i]
		clauses = append(clauses, sat.BitsAlwaysEqual(aBit, bBit)...)
	}
	fmt.Println("VectorsAreEqual", clauses)
	return clauses
}

func (constrain BitwiseLogicalConstrain) Materialize(problem *Problem) []sat.Clause {
	// returns sat.formula
	// all vectors have proper width (8)

	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	y := problem.Vectors[constrain.YIndex]

	return VectorsBitwise(a, b, y, constrain.BitConstrain)
}

func (constrain BitwiseLogicalConstrain) AddToProblem(problem *Problem) {
	problem.AddNewConstrain(constrain)
}

func (constrain BitwiseLogicalConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	return "bitwise_logical (not implemented)"
}

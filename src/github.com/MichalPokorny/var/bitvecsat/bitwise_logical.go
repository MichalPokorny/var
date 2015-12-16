package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
)

type BitConstrain func(int, int, int) []sat.Clause

type BitwiseLogicalConstrain struct {
	AIndex int
	BIndex int
	YIndex int

	BitConstrain BitConstrain
}

func (constrain BitwiseLogicalConstrain) Materialize(problem *Problem) []sat.Clause {
	// returns sat.formula
	// all vectors have proper width (8)

	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	y := problem.Vectors[constrain.YIndex]

	width := a.Width()

	if ((width != b.Width()) || (width != y.Width())) {
		panic("unequal bit widths")
	}

	clauses := make([]sat.Clause, 0) // TODO: exact width
	for i := 0; i < width; i++ {
		a := a.SatVarIndices[i]
		b := b.SatVarIndices[i]
		y := y.SatVarIndices[i]
		clauses = append(clauses, constrain.BitConstrain(a, b, y)...)
	}
	return clauses
}

func (constrain BitwiseLogicalConstrain) AddToProblem(problem *Problem) {
	problem.AddNewConstrain(constrain)
}

func OrConstrain(a int, b int, y int) []sat.Clause {
	return []sat.Clause{
		sat.NewClause(true, a, true, b, false, y),
		sat.NewClause(false, a, false, b, true, y),
	}
}

// TODO: optimize?
func AndConstrain(a int, b int, y int) []sat.Clause {
	return []sat.Clause{
		sat.NewClause(true, a, true, b, false, y),
		sat.NewClause(true, a, false, b, false, y),
		sat.NewClause(false, a, true, b, false, y),
		sat.NewClause(false, a, false, b, true, y),
	}
}

// TODO: optimize?
func XorConstrain(a int, b int, y int) []sat.Clause {
	return []sat.Clause{
		sat.NewClause(true, a, true, b, false, y),
		sat.NewClause(true, a, false, b, true, y),
		sat.NewClause(false, a, true, b, true, y),
		sat.NewClause(false, a, false, b, false, y),
	}
}

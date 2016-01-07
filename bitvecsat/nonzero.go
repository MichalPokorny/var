package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
)

type NonzeroConstrain struct {
	AIndex int
}

func (constrain NonzeroConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	clause := sat.Clause{Literals: make([]sat.Literal, width)}
	for i := uint(0); i < width; i++ {
		clause.Literals[i] = sat.Literal{
			Variable: a.SatVarIndices[i],
			Positive: true,
		}
	}
	return []sat.Clause{clause}
}

func (constrain NonzeroConstrain) AddToProblem(problem *Problem) {
	problem.AddNewConstrain(constrain)
}

func (constrain NonzeroConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	return "nonzero (todo)"
}

package bitvecsat

import "github.com/MichalPokorny/var/sat"

type LiteralConstrain struct {
	AIndex int
	Value int
}

func (constrain *LiteralConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	clauses := make([]sat.Clause, width)
	for i := uint(0); i < width; i++ {
		isHigh := (constrain.Value & (1 << i)) == (1 << i)
		clauses[i] = sat.NewClause(isHigh, a.SatVarIndices[i])
	}
	return clauses
}

func (constrain *LiteralConstrain) AddToProblem(problem *Problem) {
	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	if constrain.Value >= (1 << width) {
		panic("value too large for vector")
	}
	problem.AddNewConstrain(constrain)
}

func (constrain *LiteralConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	return "literal (not implemented)"
}

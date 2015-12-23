package bitvecsat

import "github.com/MichalPokorny/var/sat"

// Used for both division and modulo.
type DivideConstrain struct {
	AIndex int
	BIndex int

	RatioIndex int
	RemainderIndex int
}

// A / B = Ratio (with remainder Remainder) is true iff:
//
// (B > 0) && (B * Ratio + Remainder = A) && (Remainder < B)

// TODO: more effective: (B != 0)

func (constrain *DivideConstrain) Materialize(problem *Problem) []sat.Clause {
	// Fully materialized by subconstrains.
	return []sat.Clause{}
}

func (constrain *DivideConstrain) AddToProblem(problem *Problem) {
	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	width := a.Width

	// TODO: hacky!
	if constrain.RatioIndex == 0 && constrain.RemainderIndex == 0 {
		panic("please bind either ratio or remainder")
	}
	if constrain.RatioIndex == 0 {
		constrain.RatioIndex = problem.AddNewVector(width)
	}
	if constrain.RemainderIndex == 0 {
		constrain.RemainderIndex = problem.AddNewVector(width)
	}
	ratio := problem.Vectors[constrain.RatioIndex]
	remainder := problem.Vectors[constrain.RemainderIndex]
	if b.Width != width || ratio.Width != width || remainder.Width != width {
		panic("bad width")
	}

	// (B * Ratio + Remainder) = A
	bTimesRatio := problem.AddNewVector(width)

	bTimesRatioConstrain := MultiplyConstrain{
		AIndex: constrain.BIndex,
		BIndex: constrain.RatioIndex,
		ProductIndex: bTimesRatio,
	}
	bTimesRatioConstrain.AddToProblem(problem)

	plusRemainder := PlusConstrain{
		AIndex: bTimesRatio,
		BIndex: constrain.RemainderIndex,
		SumIndex: constrain.AIndex,
	}
	plusRemainder.AddToProblem(problem)

	// Remainder < B
	lt := OrderingConstrain{
		AIndex: constrain.RemainderIndex,
		BIndex: constrain.BIndex,
		Type: LT,
	}
	lt.AddToProblem(problem)

	// B != 0
	// TODO: Make this faster.
	// TODO: Make it possible to work with literals directly?
	zero := problem.AddNewVector(width)
	zeroLc := LiteralConstrain{AIndex: zero, Value: 0}
	zeroLc.AddToProblem(problem)
	gtz := OrderingConstrain{
		AIndex: zero,
		BIndex: constrain.BIndex,
		Type: LT,
	}
	gtz.AddToProblem(problem)

	problem.AddNewConstrain(constrain)
}

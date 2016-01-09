package bitvecsat

// TODO: right shift

import (
	"fmt"
	"github.com/MichalPokorny/var/sat"
)

type ShiftLeftConstrain struct {
	AIndex int
	AmountIndex int  // must be uint
	YIndex int

	MaybeShiftedIndices []int

	constantZeroIndex int
}

func (constrain *ShiftLeftConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	amount := problem.Vectors[constrain.AmountIndex]

	width := a.Width

	clauses := make([]sat.Clause, 0)

	maxShift := uint(0)
	for i := uint(0); (1 << i) <= width; i++ {
		maxShift = i
	}

	for i := uint(0); i < maxShift; i++ {
		shift := uint(1 << i)

		maybeShiftedCurrent := problem.Vectors[constrain.MaybeShiftedIndices[i]]

		var previous Vector
		if i == 0 {
			previous = a
		} else {
			previous = problem.Vectors[constrain.MaybeShiftedIndices[i - 1]]
		}

		shiftedPrevious := Vector{
			Width: width,
			SatVarIndices: make([]int, shift),
		}
		for i := 0; i < int(shift); i++ {
			shiftedPrevious.SatVarIndices[i] = constrain.constantZeroIndex
		}
		shiftedPrevious.SatVarIndices = append(
			shiftedPrevious.SatVarIndices,
			previous.SatVarIndices[0:width - shift]...)
		for j := uint(0); j < width; j++ {
			clauses = append(clauses, sat.BitIfThenElse(maybeShiftedCurrent.SatVarIndices[j], amount.SatVarIndices[i], shiftedPrevious.SatVarIndices[j], previous.SatVarIndices[j])...)
		}
	}
	return clauses
}

func (constrain *ShiftLeftConstrain) AddToProblem(problem *Problem) {
	a := problem.Vectors[constrain.AIndex]
	amount := problem.Vectors[constrain.AmountIndex]
	y := problem.Vectors[constrain.YIndex]

	width := a.Width

	if width != y.Width {
		// TODO: is this needed?
		panic("unequal widths")
	}

	// NOTE: very large shifts are ignored

	maxShift := uint(0)
	for i := uint(0); (1 << i) <= width; i++ {
		maxShift = i
	}

	if amount.Width < maxShift {
		panic(fmt.Sprintf("wrong width of amount (got %d, expected at least %d)", amount.Width, maxShift))
	}

	constrain.MaybeShiftedIndices = make([]int, maxShift + 1)

	for i := uint(0); i < maxShift; i++ {
		if i == maxShift - 1 {
			constrain.MaybeShiftedIndices[i] = constrain.YIndex
		} else {
			constrain.MaybeShiftedIndices[i] = problem.AddNewVector(width)
		}
	}

	idx := problem.AddNewVector(1)
	LiteralConstrain{AIndex: idx, Value: 0}.AddToProblem(problem)
	constrain.constantZeroIndex = problem.Vectors[idx].SatVarIndices[0]

	problem.AddNewConstrain(constrain)

}

func (constrain *ShiftLeftConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	str := fmt.Sprintf("shift_left(#%d[== %d] << #%d[== %d] = #%d[== %d])",
		constrain.AIndex,
		problem.GetValueInAssignment(assignment, constrain.AIndex),
		constrain.AmountIndex,
		problem.GetValueInAssignment(assignment, constrain.AmountIndex),
		constrain.YIndex,
		problem.GetValueInAssignment(assignment, constrain.YIndex))

	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	str += fmt.Sprintf("\n%v\n", problem.Vectors)
	str += fmt.Sprintf("\n%v\n", assignment)

	maxShift := uint(0)
	for i := uint(0); (1 << i) <= width; i++ {
		maxShift = i
	}

	for i := uint(0); i < maxShift; i++ {
		str += fmt.Sprintf("\nmaybeShiftedIndices[%d]=%d", i, problem.GetValueInAssignment(assignment, constrain.MaybeShiftedIndices[i]))
	}

	return str
}

type ShiftRightConstrain struct {
	AIndex int
	AmountIndex int
	YIndex int

	MaybeShiftedIndices []int

	constantZeroIndex int
}

func (constrain *ShiftRightConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	amount := problem.Vectors[constrain.AmountIndex]

	width := a.Width

	clauses := make([]sat.Clause, 0)

	maxShift := uint(0)
	for i := uint(0); (1 << i) <= width; i++ {
		maxShift = i
	}

	for i := uint(0); i < maxShift; i++ {
		shift := uint(1 << i)

		maybeShiftedCurrent := problem.Vectors[constrain.MaybeShiftedIndices[i]]

		var previous Vector
		if i == 0 {
			previous = a
		} else {
			previous = problem.Vectors[constrain.MaybeShiftedIndices[i - 1]]
		}

		shiftedPrevious := Vector{
			Width: width,
			SatVarIndices: previous.SatVarIndices[shift:],
		}
		for i := 0; i < int(shift); i++ {
			shiftedPrevious.SatVarIndices = append(shiftedPrevious.SatVarIndices, constrain.constantZeroIndex)
		}
		for j := uint(0); j < width; j++ {
			clauses = append(clauses, sat.BitIfThenElse(maybeShiftedCurrent.SatVarIndices[j], amount.SatVarIndices[i], shiftedPrevious.SatVarIndices[j], previous.SatVarIndices[j])...)
		}
	}
	return clauses
}

func (constrain *ShiftRightConstrain) AddToProblem(problem *Problem) {
	a := problem.Vectors[constrain.AIndex]
	amount := problem.Vectors[constrain.AmountIndex]
	y := problem.Vectors[constrain.YIndex]

	width := a.Width

	if width != y.Width {
		// TODO: is this needed?
		panic("unequal widths")
	}

	// NOTE: very large shifts are ignored

	maxShift := uint(0)
	for i := uint(0); (1 << i) <= width; i++ {
		maxShift = i
	}

	if amount.Width < maxShift {
		panic(fmt.Sprintf("wrong width of amount (got %d, expected at least %d)", amount.Width, maxShift))
	}

	constrain.MaybeShiftedIndices = make([]int, maxShift + 1)

	for i := uint(0); i < maxShift; i++ {
		if i == maxShift - 1 {
			constrain.MaybeShiftedIndices[i] = constrain.YIndex
		} else {
			constrain.MaybeShiftedIndices[i] = problem.AddNewVector(width)
		}
	}

	idx := problem.AddNewVector(1)
	LiteralConstrain{AIndex: idx, Value: 0}.AddToProblem(problem)
	constrain.constantZeroIndex = problem.Vectors[idx].SatVarIndices[0]

	problem.AddNewConstrain(constrain)

}

func (constrain *ShiftRightConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	str := fmt.Sprintf("shift_right(#%d[== %d] << #%d[== %d] = #%d[== %d])",
		constrain.AIndex,
		problem.GetValueInAssignment(assignment, constrain.AIndex),
		constrain.AmountIndex,
		problem.GetValueInAssignment(assignment, constrain.AmountIndex),
		constrain.YIndex,
		problem.GetValueInAssignment(assignment, constrain.YIndex))

	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	str += fmt.Sprintf("\n%v\n", problem.Vectors)
	str += fmt.Sprintf("\n%v\n", assignment)

	maxShift := uint(0)
	for i := uint(0); (1 << i) <= width; i++ {
		maxShift = i
	}

	for i := uint(0); i < maxShift; i++ {
		str += fmt.Sprintf("\nmaybeShiftedIndices[%d]=%d", i, problem.GetValueInAssignment(assignment, constrain.MaybeShiftedIndices[i]))
	}

	return str
}

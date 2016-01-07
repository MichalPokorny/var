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
	ShiftedIndices []int
}

func (constrain *ShiftLeftConstrain) Materialize(problem *Problem) []sat.Clause {
	// TODO: barrel shifter

	a := problem.Vectors[constrain.AIndex]
	amount := problem.Vectors[constrain.AmountIndex]

	width := a.Width

	clauses := make([]sat.Clause, 0)
	for i := uint(0); (1 << i) < width; i++ {
		shift := uint(1 << i)

		shiftedPrevious := problem.Vectors[constrain.ShiftedIndices[i]]
		maybeShiftedCurrent := problem.Vectors[constrain.MaybeShiftedIndices[i]]

		var previous Vector
		if i == 0 {
			previous = a
		} else {
			previous = problem.Vectors[constrain.MaybeShiftedIndices[i - 1]]
		}

		for j := uint(0); j < width; j++ {
			if j < shift {
				clauses = append(clauses, sat.BitIsFalse(shiftedPrevious.SatVarIndices[j])...)
			} else {
				clauses = append(clauses, sat.BitsAlwaysEqual(shiftedPrevious.SatVarIndices[j], previous.SatVarIndices[j - shift])...)
			}

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

	// TODO: check width of 'amount'!
	if maxShift != amount.Width {
		fmt.Println("got", amount.Width, "expected", maxShift)
		panic("wrong width of amount")
	}

	constrain.ShiftedIndices = make([]int, maxShift + 1)
	constrain.MaybeShiftedIndices = make([]int, maxShift + 1)

	for i := uint(0); (1 << i) <= width; i++ {
		constrain.ShiftedIndices[i] = problem.AddNewVector(width)
		if i == maxShift {
			constrain.MaybeShiftedIndices[i] = constrain.YIndex
		} else {
			constrain.MaybeShiftedIndices[i] = problem.AddNewVector(width)
		}
	}


	// TODO: barrel shifter

	problem.AddNewConstrain(constrain)
}

func (constrain *ShiftLeftConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	return "shift left (not implemented)"
}

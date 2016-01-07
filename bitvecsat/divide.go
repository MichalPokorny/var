package bitvecsat

import (
	"fmt"
	"strconv"
	"github.com/MichalPokorny/var/sat"
)

// Used for both division and modulo.
type DivideConstrain struct {
	AIndex int
	BIndex int
	// assuming length(A) == length(B), marked as 'k'

	RatioIndex int
	RemainderIndex int

	// k subresults, first one is padded with k zeros on left
	// each subresult has length k+1
	subresultIndices []int

	startMinusDividendIndices []int
	startAtIndices []int

	extendedDividendIndex int

	lteConstrains []OrderingConstrain
}

// A / B = Ratio (with remainder Remainder) is true iff:
//
// (B > 0) && (B * Ratio + Remainder = A) && (Remainder < B)

// TODO: more effective: (B != 0)

func (constrain *DivideConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	ratio := problem.Vectors[constrain.RatioIndex]

	clauses := make([]sat.Clause, 0)
	// (if startAt > ([0] . divisor) then stopAt := startMinusDividend else stopAt := startAt)
	for i := 0; i < int(width); i++ {
		startAt := problem.Vectors[constrain.startAtIndices[i]]
		stopAt := problem.Vectors[constrain.subresultIndices[i]]
		resultBit := ratio.SatVarIndices[int(width) - 1 - i]
		startMinusDividend := problem.Vectors[constrain.startMinusDividendIndices[i]]

		for j := 0; j < int(width); j++ {
			subresultBit := stopAt.SatVarIndices[j]
			clause := sat.BitIfThenElse(subresultBit, resultBit, startMinusDividend.SatVarIndices[j], startAt.SatVarIndices[j])
			clauses = append(clauses, clause...)
		}
	}
	return clauses
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

	// pad A with k zeros
	zerosIndex := problem.AddNewVector(width)
	zeros := problem.Vectors[zerosIndex]
	zerosConstrain := LiteralConstrain{AIndex: zerosIndex, Value: 0}
	zerosConstrain.AddToProblem(problem)

	constrain.subresultIndices = make([]int, width)
	for i := 0; i < int(width) - 1; i++ {
		constrain.subresultIndices[i] = problem.AddNewVector(width)
	}
	constrain.subresultIndices[width - 1] = constrain.RemainderIndex

	constrain.extendedDividendIndex = problem.AddBoundVector(width + 1, append(b.SatVarIndices, zeros.SatVarIndices[0]))

	constrain.startAtIndices = make([]int, width)
	constrain.startMinusDividendIndices = make([]int, width)
	constrain.lteConstrains = make([]OrderingConstrain, width)

	for i := 0; i < int(width); i++ {
		var startAtComposition []int
		addedFromA := []int{a.SatVarIndices[int(width) - 1 - i]}
		if i == 0 {
			startAtComposition = append(addedFromA, zeros.SatVarIndices...)
		} else {
			startAtComposition = append(addedFromA, problem.Vectors[constrain.subresultIndices[i - 1]].SatVarIndices...)
		}
		startAtIndex := problem.AddBoundVector(width + 1, startAtComposition)
		constrain.startAtIndices[i] = startAtIndex

		startMinusDividend := problem.AddNewVector(width + 1)
		constrain.startMinusDividendIndices[i] = startMinusDividend
		plusConstrain := PlusConstrain{
			AIndex: constrain.extendedDividendIndex,
			BIndex: startMinusDividend,
			SumIndex: startAtIndex,
		}
		plusConstrain.AddToProblem(problem)

		resultBit := ratio.SatVarIndices[int(width) - 1 - i]
		constrain.lteConstrains[i] = OrderingConstrain{
			AIndex: constrain.extendedDividendIndex,
			BIndex: startAtIndex,
			Type: LTE,
			IsQuery: true,
			QueryResultBitIndex: resultBit,
		}
		constrain.lteConstrains[i].AddToProblem(problem)
	}

	NonzeroConstrain{AIndex: constrain.BIndex}.AddToProblem(problem)

	problem.AddNewConstrain(constrain)
}

func (constrain *DivideConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	str := "divide(#" + strconv.Itoa(constrain.AIndex) + " / #" + strconv.Itoa(constrain.BIndex) + " = #" + strconv.Itoa(constrain.RatioIndex) + " (rem #" + strconv.Itoa(constrain.RemainderIndex) + "))"
	str += "\ndivisor=" + strconv.Itoa(problem.GetValueInAssignment(assignment, constrain.AIndex))
	str += "\ndividend=" + strconv.Itoa(problem.GetValueInAssignment(assignment, constrain.BIndex))
	str += "\nratio=" + strconv.Itoa(problem.GetValueInAssignment(assignment, constrain.RatioIndex))
	str += "\nremainder=" + strconv.Itoa(problem.GetValueInAssignment(assignment, constrain.RemainderIndex))
	str += "\n"
	str += "\nextendedDividend=" + strconv.Itoa(problem.GetValueInAssignment(assignment, constrain.extendedDividendIndex))
	str += "\n"

	a := problem.Vectors[constrain.AIndex]
	width := a.Width
	for i := 0; i < int(width); i++ {
		str += fmt.Sprintf("\nstartAt[%d]=#%d[== %d]", i, constrain.startAtIndices[i], problem.GetValueInAssignment(assignment, constrain.startAtIndices[i]))
		str += fmt.Sprintf("\nstartMinusDividend[%d]=%d", i, problem.GetValueInAssignment(assignment, constrain.startMinusDividendIndices[i]))
		bit := assignment[problem.Vectors[constrain.RatioIndex].SatVarIndices[int(width) - 1 - i]]
		str += fmt.Sprintf("\nresultBit=%v", bit)
		str += fmt.Sprintf("\nsubresults[%d]=%d", i, problem.GetValueInAssignment(assignment, constrain.subresultIndices[i]))
		str += "\n" + constrain.lteConstrains[i].Dump(problem, assignment)
		str += "\n"
	}
	return str
}

package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
)

type MultiplyConstrain struct {
	AIndex int
	BIndex int
	ProductIndex int

	// [i] is the last (width-i) bits of product of A with B[i]
	SubresultIndices []int

	SubsumIndices []int
}

func (constrain *MultiplyConstrain) Materialize(problem *Problem) []sat.Clause {
	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]

	width := a.Width

	clauses := make([]sat.Clause, 0)

	// O(N^2) grammar school multiplication
	// TODO: implement O(N log N) as well? compare the approaches?
	for i := 0; i < int(width); i++ {
		subresult := problem.Vectors[constrain.SubresultIndices[i]]
		for j := 0; j < i; j++ {
			clauses = append(clauses, sat.BitIsFalse(subresult.SatVarIndices[j])...)
		}
		for j := 0; j < (int(width) - i); j++ {
			clauses = append(clauses, sat.AndConstrain(a.SatVarIndices[j], b.SatVarIndices[i], subresult.SatVarIndices[i + j])...)
		}
	}

	return clauses
}

// TODO: O(N log N) addition, not O(N^2)
func (constrain *MultiplyConstrain) AddToProblem(problem *Problem) {
	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	product := problem.Vectors[constrain.ProductIndex]

	width := a.Width

	if (width != b.Width) || (width != product.Width) {
		// TODO: is this needed?
		panic("unequal widths")
	}

	constrain.SubresultIndices = make([]int, width)
	constrain.SubsumIndices = make([]int, width)

	for i := 0; i < int(width); i++ {
		constrain.SubresultIndices[i] = problem.AddNewVector(width)
	}

	// Subsum[0] = Subresult[0]
	constrain.SubsumIndices[0] = constrain.SubresultIndices[0]
	// Subsum[width - 1] = product
	constrain.SubsumIndices[width - 1] = constrain.ProductIndex

	if width == 1 {
		// Extra: also need to set SubresultIndices[0] to product.
		constrain.SubresultIndices[0] = constrain.ProductIndex
	}

	for i := 1; i < int(width - 1); i++ {
		constrain.SubsumIndices[i] = problem.AddNewVector(width)
	}

	for i := 0; i < int(width) - 1; i++ {
		plusConstrain := PlusConstrain{
			AIndex: constrain.SubresultIndices[i + 1],
			BIndex: constrain.SubsumIndices[i],
			SumIndex: constrain.SubsumIndices[i + 1],
		}
		plusConstrain.AddToProblem(problem)
	}

	problem.AddNewConstrain(constrain)
}

func (constrain *MultiplyConstrain) Dump(problem *Problem, assignment sat.Assignment) string {
	return "multiply (not implemented)"
}

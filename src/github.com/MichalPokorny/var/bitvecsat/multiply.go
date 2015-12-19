package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
	"fmt"
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
	product := problem.Vectors[constrain.ProductIndex]

	width := a.Width

	clauses := make([]sat.Clause, 0)

	// O(N^2) grammar school multiplication
	// TODO: implement O(N log N) as well? compare the approaches?
	for i := 0; i < int(width); i++ {
		subresult := problem.Vectors[constrain.SubresultIndices[i]]
		fmt.Println("subresult[", i, "]=", subresult)
		for j := 0; j < i; j++ {
			clauses = append(clauses, sat.BitIsFalse(subresult.SatVarIndices[j])...)
		}
		for j := 0; j < (int(width) - i); j++ {
			clauses = append(clauses, sat.AndConstrain(a.SatVarIndices[j], b.SatVarIndices[i], subresult.SatVarIndices[i + j])...)
		}
	}
	fmt.Println("constrs:", clauses)

	// TODO: faster addition
	for i := 0; i < int(width); i++ {
		subsum := problem.Vectors[constrain.SubsumIndices[i]]
		fmt.Println("subsum[", i, "]=", subsum)
		if i == 0 {
			clauses = append(clauses, VectorsAreEqual(subsum, problem.Vectors[constrain.SubresultIndices[0]])...)
		}
		if i == int(width) - 1 {
			clauses = append(clauses, VectorsAreEqual(subsum, product)...)
		}

		fmt.Println("constrs", i, ":", clauses)
	}

	return clauses
}

// TODO: this approach probably creates holes... :(
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

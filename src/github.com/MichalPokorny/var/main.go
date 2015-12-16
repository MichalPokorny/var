package main

import (
	"strconv"
	"fmt"
	"github.com/MichalPokorny/var/sat"
	"github.com/MichalPokorny/var/bitvecsat"
)

func ExhaustAllSolutions(formula sat.Formula) {
	for {
		solution := sat.Solve(formula)
		if solution == nil {
			fmt.Println("No more solutions.")
			break
		}
		fmt.Println(solution.String())
		formula.Clauses = append(formula.Clauses, solution.MakeForbiddingClause())
	}
}

func FindOneSolution(formula sat.Formula) {
	fmt.Println(sat.Solve(formula).String())
}

/*
func TestBitvecsatPrint() {
	a := bitvecsat.Vector{SatVarIndices: []int{0, 1, 2, 3, 4, 5, 6, 7}}
	b := bitvecsat.Vector{SatVarIndices: []int{8, 9, 10, 11, 12, 13, 14, 15}}
	c := bitvecsat.Vector{SatVarIndices: []int{16, 17, 18, 19, 20, 21, 22, 23}}
	problem := bitvecsat.Problem{Vectors: []bitvecsat.Vector{a, b, c}}
	fmt.Println(problem)
}

func TestBitCarry() {
	formula := sat.Formula{
		Clauses: bitvecsat.AddBitCarryClause(0, 1, 2, 3),
	}
	ExhaustAllSolutions(formula)
}
*/

func TestAddition() {
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector()
	b := problem.AddNewVector()
	c := problem.AddNewVector()

	plus_constrain := bitvecsat.PlusConstrain{AIndex: a, BIndex: b, SumIndex: c}
	plus_constrain.AddToProblem(&problem)

	problem.PrepareSat(8)
	fmt.Println(problem)

	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := sat.Solve(formula)
		if solution == nil {
			fmt.Println("No more solutions.")
			break
		}
		//fmt.Println(solution.String())

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)
		cValue := problem.GetValueInAssignment(solution, c)

		fmt.Println("A=" + strconv.Itoa(aValue) + " B=" + strconv.Itoa(bValue) + " C=" + strconv.Itoa(cValue));
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}
}

func main() {
	TestAddition()
}

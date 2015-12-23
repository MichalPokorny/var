package main

import (
	"strconv"
	"fmt"
	"github.com/MichalPokorny/var/sat"
	"github.com/MichalPokorny/var/sat/dfs"
	"github.com/MichalPokorny/var/sat/cdcl"
	"github.com/MichalPokorny/var/sat/dpll"
	"github.com/MichalPokorny/var/bitvecsat"
)

func solveFormula(formula sat.Formula) sat.Assignment {
	if false {
		// DPLL
		return dpll.Solve(formula, sat.MakeEmptyAssignment(formula))
	}

	if false {
		// DFS
		return dfs.Solve(formula)
	}

	return cdcl.Solve(formula)
}

func ExhaustAllSolutions(formula sat.Formula) {
	for {
		solution := solveFormula(formula)
		if solution == nil {
			fmt.Println("No more solutions.")
			break
		}
		fmt.Println(solution.String())
		formula.Clauses = append(formula.Clauses, solution.MakeForbiddingClause())
	}
}

func FindOneSolution(formula sat.Formula) {
	fmt.Println(solveFormula(formula).String())
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

func ShowAddition() {
	width := uint(2)
	problem := bitvecsat.Problem{}
	a := problem.AddNewVector(width)
	b := problem.AddNewVector(width)
	c := problem.AddNewVector(width)

	multiply_constrain := &bitvecsat.MultiplyConstrain{AIndex: a, BIndex: b, ProductIndex: c}
	multiply_constrain.AddToProblem(&problem)
	//shift_constrain := &bitvecsat.ShiftLeftConstrain{AIndex: a, AmountIndex: b, YIndex: c}
	//shift_constrain.AddToProblem(&problem)

	//lte_constrain := bitvecsat.OrderingConstrain{AIndex: a, BIndex: b}
	//lte_constrain.AddToProblem(&problem)

	problem.PrepareSat()

	fmt.Println(problem)

	formula := problem.MakeSatFormula()
	fmt.Println("formula: " + formula.String())
	forbidders := make([]sat.Clause, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		//fmt.Println(formula.Clauses)
		solution := solveFormula(formula)
		if solution == nil {
			fmt.Println("No more solutions.")
			break
		}
		//fmt.Println(solution.String())


		/*
		fmt.Println()
		fmt.Println(solution, " len=", len(solution))
		fmt.Println("constrains:")
		for _, constrain := range(problem.Constrains) {
			fmt.Println(constrain)
		}
		fmt.Println("vectors:")
		for i, vector := range(problem.Vectors) {
			fmt.Println("[", i, "]=", vector)
		}
		*/

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)
		cValue := problem.GetValueInAssignment(solution, c)

		aString := problem.GetBitsInAssignment(solution, a)
		bString := problem.GetBitsInAssignment(solution, b)
		cString := problem.GetBitsInAssignment(solution, c)

		/*
		fmt.Println("shifted");
		for i := 0; i < len(shift_constrain.ShiftedIndices); i++ {
			value := problem.GetValueInAssignment(solution, shift_constrain.ShiftedIndices[i])
			bits := problem.GetBitsInAssignment(solution, shift_constrain.ShiftedIndices[i])
			fmt.Println("[", i, "]=" + strconv.Itoa(value) + "=" + bits);
		}
		fmt.Println("maybe-shifted");
		for i := 0; i < len(shift_constrain.MaybeShiftedIndices); i++ {
			value := problem.GetValueInAssignment(solution, shift_constrain.MaybeShiftedIndices[i])
			bits := problem.GetBitsInAssignment(solution, shift_constrain.MaybeShiftedIndices[i])
			fmt.Println("[", i, "]=" + strconv.Itoa(value) + "=" + bits);
		}
		*/

		fmt.Println("A=" + strconv.Itoa(aValue) + "=" + aString + " B=" + strconv.Itoa(bValue) + "=" + bString + " C=" + strconv.Itoa(cValue) + "=" + cString);
		// TODO: fix this!
		//fmt.Println("A=" + strconv.Itoa(aValue) + "=" + aString + " B=" + strconv.Itoa(bValue) + "=" + bString);
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}
}

func ShowSat() {
	formula := sat.Formula{
		Clauses: []sat.Clause{
			// example from wiki
			// https://en.wikipedia.org/wiki/Conflict-Driven_Clause_Learning
			sat.NewClause(false, 0, true, 1, true, 2),
			sat.NewClause(true, 0, true, 2, true, 3),
			sat.NewClause(true, 0, true, 2, false, 3),
			sat.NewClause(true, 0, false, 2, true, 3),
			sat.NewClause(true, 0, false, 2, false, 3),
			sat.NewClause(false, 1, false, 2, true, 3),
			sat.NewClause(false, 0, true, 1, false, 2),
			sat.NewClause(false, 0, false, 1, true, 2),
		},
	}
	assignment := cdcl.Solve(formula)
	fmt.Println(assignment)
}

func main() {
	cdcl.Init()

	ShowAddition()
//	ShowSat()
}

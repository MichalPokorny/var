package main

import (
	"strconv"
	//"fmt"
	"sort"
	"testing"
	"github.com/MichalPokorny/var/sat"
	"github.com/MichalPokorny/var/sat/dfs"
	"github.com/MichalPokorny/var/sat/dpll"
	"github.com/MichalPokorny/var/sat/cdcl"
	"github.com/MichalPokorny/var/bitvecsat"
)

func solve(formula sat.Formula) sat.Assignment {
	if false {
		// DPLL
		return dpll.Solve(formula, sat.MakeEmptyAssignment(formula))
	}

	if false {
		// DFS
		return dfs.Solve(formula)
	}

	// CDCL
	return cdcl.Solve(formula)
}

type resultSorter struct {
	Elements [][]int
}

func (sorter resultSorter) Len() int {
	return len(sorter.Elements)
}

func (sorter resultSorter) Less(i, j int) bool {
	for idx := 0; idx < len(sorter.Elements[i]); idx++ {
		if sorter.Elements[i][idx] < sorter.Elements[j][idx] {
			return true
		}
		if sorter.Elements[i][idx] > sorter.Elements[j][idx] {
			return false
		}
	}
	return false
}

func (sorter resultSorter) Swap(i, j int) {
	tmp := sorter.Elements[i]
	sorter.Elements[i] = sorter.Elements[j]
	sorter.Elements[j] = tmp
}

func intSlicesEqual(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func resultSetsEqual(a [][]int, b [][]int) bool {
	sort.Sort(resultSorter{Elements: a})
	sort.Sort(resultSorter{Elements: b})
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !intSlicesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func resultSetContains(results [][]int, result []int) bool {
	for i := 0; i < len(results); i++ {
		if intSlicesEqual(results[i], result) {
			return true
		}
	}
	return false
}

func getResultsDifference(a [][]int, b [][]int) [][]int {
	difference := make([][]int, 0)
	for _, slice := range(a) {
		if !resultSetContains(b, slice) {
			difference = append(difference, slice)
		}
	}
	return difference
}

type binaryOperator func(a, b int, width uint) int;

func testBinaryOperator(t *testing.T, width uint, a int, b int, c int, problem bitvecsat.Problem, operator binaryOperator) {
	problem.PrepareSat()
	formula := problem.MakeSatFormula()
	t.Logf("Formula: %v", formula)
	forbidders := make([]sat.Clause, 0)

	// Contains 3-tuples of (A, B, C)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := solve(formula)
		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)
		cValue := problem.GetValueInAssignment(solution, c)

		if true {
			aString := problem.GetBitsInAssignment(solution, a)
			bString := problem.GetBitsInAssignment(solution, b)
			cString := problem.GetBitsInAssignment(solution, c)

			t.Log("A=" + strconv.Itoa(aValue) + "=" + aString + " B=" + strconv.Itoa(bValue) + "=" + bString + " C=" + strconv.Itoa(cValue) + "=" + cString);
		}

		foundSolutions = append(foundSolutions, []int{aValue, bValue, cValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		for b := 0; b < (1 << width); b++ {
			c := operator(a, b, width)
			expectedSolutions = append(expectedSolutions, []int{a, b, c})
		}
	}

	if !resultSetsEqual(foundSolutions, expectedSolutions) {
		t.Errorf("unexpected solutions: found %v, expected %v", foundSolutions, expectedSolutions)
		t.Errorf("extra: %v", getResultsDifference(foundSolutions, expectedSolutions))
		t.Errorf("missing: %v", getResultsDifference(expectedSolutions, foundSolutions))
	}
}

func operatorPlus(a int, b int, width uint) int {
	return (a + b) % (1 << width);
}

func TestAddition(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		plus_constrain := bitvecsat.PlusConstrain{AIndex: a, BIndex: b, SumIndex: c}
		plus_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorPlus)
	}
}

func operatorOr(a int, b int, width uint) int {
	return a | b;
}

func TestOr(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		or_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.OrConstrain}
		or_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorOr)
	}
}

func operatorAnd(a int, b int, width uint) int {
	return a & b;
}

func TestAnd(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		and_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.AndConstrain}
		and_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorAnd)
	}
}

func operatorXor(a int, b int, width uint) int {
	return a ^ b;
}

func TestXor(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		xor_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.XorConstrain}
		xor_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorXor)
	}
}

func operatorEquiv(a int, b int, width uint) int {
	return (1 << width - 1) - (a ^ b);
}

func TestEquiv(t *testing.T) {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		equiv_constrain := bitvecsat.BitwiseLogicalConstrain{AIndex: a, BIndex: b, YIndex: c, BitConstrain: sat.EquivConstrain}
		equiv_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorEquiv)
	}
}

type binaryRelation func(a, b int, width uint) bool;

func testBinaryRelation(t *testing.T, width uint, a int, b int, problem bitvecsat.Problem, relation binaryRelation) {
	problem.PrepareSat()
	formula := problem.MakeSatFormula()
	forbidders := make([]sat.Clause, 0)

	// t.Log(problem)
	// t.Log(formula)

	// Contains 3-tuples of (A, B, C)
	foundSolutions := make([][]int, 0)

	for {
		formula.Clauses = append(formula.Clauses, forbidders...)
		solution := solve(formula)
		if solution == nil {
			break
		}

		aValue := problem.GetValueInAssignment(solution, a)
		bValue := problem.GetValueInAssignment(solution, b)

		foundSolutions = append(foundSolutions, []int{aValue, bValue})
		forbidders = append(forbidders, solution.MakeForbiddingClause())
	}

	expectedSolutions := make([][]int, 0)
	for a := 0; a < (1 << width); a++ {
		for b := 0; b < (1 << width); b++ {
			if relation(a, b, width) {
				expectedSolutions = append(expectedSolutions, []int{a, b})
			}
		}
	}

	if !resultSetsEqual(foundSolutions, expectedSolutions) {
		t.Errorf("unexpected solutions: found %v, expected %v", foundSolutions, expectedSolutions)
	}
}

func relationLte(a int, b int, width uint) bool {
	return a <= b;
}

func TestLte(t *testing.T) {
	for width := uint(1); width < 4; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)

		lte_constrain := bitvecsat.OrderingConstrain{AIndex: a, BIndex: b, Type: bitvecsat.LTE}
		lte_constrain.AddToProblem(&problem)

		testBinaryRelation(t, width, a, b, problem, relationLte)
	}
}

func relationLt(a int, b int, width uint) bool {
	return a < b;
}

func TestLt(t *testing.T) {
	for width := uint(1); width < 4; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)

		lte_constrain := bitvecsat.OrderingConstrain{AIndex: a, BIndex: b, Type: bitvecsat.LT}
		lte_constrain.AddToProblem(&problem)

		testBinaryRelation(t, width, a, b, problem, relationLt)
	}
}

func operatorMultiply(a int, b int, width uint) int {
	return (a * b) % (1 << width);
}

func TestMultiplication(t *testing.T) {
//	for width := uint(1); width <= 4; width++ {
	for width := uint(1); width <= 3; width++ {
		problem := bitvecsat.Problem{}
		a := problem.AddNewVector(width)
		b := problem.AddNewVector(width)
		c := problem.AddNewVector(width)

		mutliply_constrain := bitvecsat.MultiplyConstrain{AIndex: a, BIndex: b, ProductIndex: c}
		mutliply_constrain.AddToProblem(&problem)

		testBinaryOperator(t, width, a, b, c, problem, operatorMultiply)
	}
}

/*
func operatorShiftLeft(a int, b int, width uint) int {
	return (a << uint(b)) % (1 << width);
}

func TestShiftLeft(t *testing.T) {
	shift := uint(1)
	width := uint(1 << shift)

	problem := bitvecsat.Problem{}
	a := problem.AddNewVector(width)
	b := problem.AddNewVector(shift)
	c := problem.AddNewVector(width)

	// TODO: maybe get this to work on bigger widths as well, instead?
	shift_constrain := bitvecsat.ShiftLeftConstrain{AIndex: a, AmountIndex: b, YIndex: c}
	shift_constrain.AddToProblem(&problem)

	testBinaryOperator(t, width, a, b, c, problem, operatorShiftLeft)
}
*/
